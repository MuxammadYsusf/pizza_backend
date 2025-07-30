package db

import (
	"database/sql"
	"fmt"
	"log"

	"github/http/copy/task4/config"

	_ "github.com/lib/pq"
)

func Postgres(cfg *config.Config) (*sql.DB, error) {
	connect := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
		cfg.SSLMode,
	)

	conn, err := sql.Open("postgres", connect)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return conn, nil
}
