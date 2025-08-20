package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github/http/copy/task4/config"
	"github/http/copy/task4/internal/storage/postgres"
	sqlc "github/http/copy/task4/internal/storage/postgres/sqlc/generated"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saidamir98/udevs_pkg/logger"
	"go.uber.org/zap"
)

type Storage interface {
	SQLC() sqlc.Querier
	Postgres() postgres.NewPostgresI
	Close()
	InTx(ctx context.Context, fn func(q sqlc.Querier) error) error
}

type store struct {
	postgresI postgres.NewPostgresI
	sqlcStore *SQLStore
}

type SQLStore struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewStore(ctx context.Context, dsn string, log logger.LoggerI) (*SQLStore, error) {
	log.Info("connecting to postgres...",
		zap.String("dsn", maskDSN(dsn)),
	)

	if dsn == "" {
		return nil, fmt.Errorf("empty dsn (BUG: must use cfg.PostgresURLFinal)")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error("failed to parse postgres config", zap.Error(err))
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Error("failed to create pool", zap.Error(err))
		return nil, fmt.Errorf("new pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = pool.Ping(pingCtx); err != nil {
		pool.Close()
		log.Error("failed to ping postgres", zap.Error(err))
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &SQLStore{
		pool:    pool,
		queries: sqlc.New(pool),
	}, nil
}

func (s *SQLStore) Queries() *sqlc.Queries { return s.queries }

func (s *SQLStore) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	q := sqlc.New(tx)
	if err = fn(q); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %w (rollback err: %v)", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

// New инициализирует общий Storage
func New(ctx context.Context, cfg *config.Config, log logger.LoggerI) (Storage, error) {
	dsn := cfg.PostgresURLFinal
	if dsn == "" {
		return nil, fmt.Errorf("empty PostgresURLFinal (BUG)")
	}

	sqlStore, err := NewStore(ctx, dsn, log)
	if err != nil {
		return nil, err
	}

	// Если нужен старый слой через database/sql:
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}
	// Можно добавить ping:
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("sql ping: %w", err)
	}

	return &store{
		postgresI: postgres.NewPostgres(db),
		sqlcStore: sqlStore,
	}, nil
}

func (s *store) SQLC() sqlc.Querier              { return s.sqlcStore.Queries() }
func (s *store) Postgres() postgres.NewPostgresI { return s.postgresI }
func (s *store) Close() {
	if s.sqlcStore != nil {
		s.sqlcStore.Close()
	}
	// Добавь закрытие db внутри postgres.NewPostgresI реализации, если нужно.
}

func (s *store) InTx(ctx context.Context, fn func(q sqlc.Querier) error) error {
	if s.sqlcStore == nil {
		return errors.New("postgres storage not initialized")
	}
	return s.sqlcStore.execTx(ctx, func(q *sqlc.Queries) error {
		return fn(q)
	})
}

// maskDSN простая маска пароля
func maskDSN(dsn string) string {
	if dsn == "" {
		return dsn
	}
	// очень упрощённо (можно использовать cfg.MaskedPostgresURL())
	i := 0
	if i = len("postgres://"); len(dsn) > i && dsn[:i] == "postgres://" {
		rest := dsn[i:]
		at := -1
		for j := 0; j < len(rest); j++ {
			if rest[j] == '@' {
				at = j
				break
			}
		}
		if at != -1 {
			cred := rest[:at]
			colon := -1
			for j := 0; j < len(cred); j++ {
				if cred[j] == ':' {
					colon = j
					break
				}
			}
			if colon != -1 {
				return "postgres://" + cred[:colon] + ":***" + "@" + rest[at+1:]
			}
		}
	}
	return dsn
}