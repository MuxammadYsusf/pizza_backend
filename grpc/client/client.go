package client

import (
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/generated/session"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceManager interface {
	Login() session.AuthServiceClient
	Pizza() pizza.PizzaServiceClient
}

type GRPCClient struct {
	login session.AuthServiceClient
	pizza pizza.PizzaServiceClient
}

func (g *GRPCClient) Login() session.AuthServiceClient {
	return g.login
}

func (g *GRPCClient) Pizza() pizza.PizzaServiceClient {
	return g.pizza
}

func NewGRPCClient() (ServiceManager, error) {
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		login: session.NewAuthServiceClient(conn),
		pizza: pizza.NewPizzaServiceClient(conn),
	}, err
}
