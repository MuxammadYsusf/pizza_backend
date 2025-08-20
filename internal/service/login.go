package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github/http/copy/task4/config"
	"github/http/copy/task4/internal/storage"
	sqlc "github/http/copy/task4/internal/storage/postgres/sqlc/generated"
	"github/http/copy/task4/pkg/util"

	pb "github/http/copy/task4/genproto"
	client "github/http/copy/task4/internal/transport/grpc/client"

	"github.com/saidamir98/udevs_pkg/logger"
)

type authService struct {
	cfg      config.Config
	log      logger.LoggerI
	storage  storage.Storage
	services client.ServiceManager
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(cfg config.Config, log logger.LoggerI, storage storage.Storage, services client.ServiceManager) *authService {
	return &authService{
		cfg:      cfg,
		log:      log,
		storage:  storage,
		services: services,
	}
}

func (s *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 1. Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("invalid name or password or email address")
	}
	if strings.ContainsAny(req.Username, "!@#$%&*()№:;<>?'\";--") {
		return nil, fmt.Errorf("invalid name")
	}
	if strings.ContainsAny(req.Password, "'\";--") {
		return nil, fmt.Errorf("invalid password")
	}
	if len(req.Username) > 15 || len(req.Password) > 50 {
		return nil, fmt.Errorf("invalid name or password")
	}

	// 2. Role assignment
	if req.Password == "admin" {
		req.Role = "admin"
	} else {
		req.Role = "user"
	}

	// 3. Транзакция через InTx
	err := s.storage.InTx(ctx, func(q sqlc.Querier) error {
		// 4. Check duplicates
		exists, err := q.IsNameTaken(ctx, req.Username)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("username already exists")
		}

		exists, err = q.IsEmailTaken(ctx, req.Email)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("email already exists")
		}

		// 5. Hash password
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			return err
		}

		// 6. Insert user
		_, err = q.CreateUser(ctx, sqlc.CreateUserParams{
			Name:     req.Username,
			Password: hashedPassword,
			Email:    req.Email,
			Role:     req.Role,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Message: "success"}, nil
}


func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("invalid name or password")
	}
	if strings.ContainsAny(req.Username, "!@#$%&*()№:;<>?") {
		return nil, fmt.Errorf("invalid name or password")
	}
	if strings.ContainsAny(req.Password, "'\";--") {
		return nil, fmt.Errorf("invalid name or password")
	}
	if len(req.Username) > 15 || len(req.Password) > 50 {
		return nil, fmt.Errorf("invalid name or password")
	}

	user, err := s.storage.SQLC().GetUserByName(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid name or password")
		}
		return nil, err
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid name or password")
	}

	return &pb.LoginResponse{
		Message: "success",
		UserId:  int32(user.ID),
		Role:    user.Role,
	}, nil
}
