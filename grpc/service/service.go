package service

import (
	"database/sql"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/generated/session"
	"github/http/copy/task4/grpc/client"
	"github/http/copy/task4/postgres"
)

type PizzaService struct {
	pizzaPostgres postgres.NewPostgresI
	service       client.ServiceManager
	pizza.UnimplementedPizzaServiceServer
}

type LoginService struct {
	loginPostgres postgres.NewPostgresI
	service       client.ServiceManager
	session.UnimplementedAuthServiceServer
}

func NewPizzaService(db *sql.DB, service client.ServiceManager) *PizzaService {
	return &PizzaService{
		pizzaPostgres: postgres.NewPostgres(db),
		service:       service,
	}
}

func NewLoginService(db *sql.DB, service client.ServiceManager) *LoginService {
	return &LoginService{
		loginPostgres: postgres.NewPostgres(db),
		service:       service,
	}
}
