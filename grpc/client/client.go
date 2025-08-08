package client

import (
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/generated/session"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceManager interface {
	Auth() session.AuthServiceClient
	Pizza() pizza.PizzaServiceClient
	Cart() pizza.CartServiceClient
	Order() pizza.OrderServiceClient
}

type GRPCClient struct {
	auth  session.AuthServiceClient
	pizza pizza.PizzaServiceClient
	cart  pizza.CartServiceClient
	order pizza.OrderServiceClient
}

func (g *GRPCClient) Auth() session.AuthServiceClient {
	return g.auth
}

func (g *GRPCClient) Pizza() pizza.PizzaServiceClient {
	return g.pizza
}

func (g *GRPCClient) Cart() pizza.CartServiceClient {
	return g.cart
}

func (g *GRPCClient) Order() pizza.OrderServiceClient {
	return g.order
}

func NewGRPCClient() (ServiceManager, error) {
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		auth:  session.NewAuthServiceClient(conn),
		pizza: pizza.NewPizzaServiceClient(conn),
		cart:  pizza.NewCartServiceClient(conn),
		order: pizza.NewOrderServiceClient(conn),
	}, err
}
