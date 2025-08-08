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
	pizza.UnimplementedCartServiceServer
	pizza.UnimplementedOrderServiceServer
}

type AuthService struct {
	loginPostgres postgres.NewPostgresI
	service       client.ServiceManager
	session.UnimplementedAuthServiceServer
}

func NewAuthnService(db *sql.DB, service client.ServiceManager) *AuthService {
	return &AuthService{
		loginPostgres: postgres.NewPostgres(db),
		service:       service,
	}
}

func NewPizzaService(db *sql.DB, service client.ServiceManager) *PizzaService {
	return &PizzaService{
		pizzaPostgres: postgres.NewPostgres(db),
		service:       service,
	}
}

func NewCartService(db *sql.DB, service client.ServiceManager) *PizzaService {
	return &PizzaService{
		pizzaPostgres: postgres.NewPostgres(db),
		service:       service,
	}
}

func NewOrderService(db *sql.DB, service client.ServiceManager) *PizzaService {
	return &PizzaService{
		pizzaPostgres: postgres.NewPostgres(db),
		service:       service,
	}
}
