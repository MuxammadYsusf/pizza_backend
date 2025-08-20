package grpc_client

import (
	pb "github/http/copy/task4/genproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceManager interface {
	Auth() pb.AuthServiceClient
	Pizza() pb.PizzaServiceClient
	Cart() pb.CartServiceClient
	Order() pb.OrderServiceClient
}

type GRPCClient struct {
	auth  pb.AuthServiceClient
	pizza pb.PizzaServiceClient
	cart  pb.CartServiceClient
	order pb.OrderServiceClient
}

func (g *GRPCClient) Auth() pb.AuthServiceClient {
	return g.auth
}

func (g *GRPCClient) Pizza() pb.PizzaServiceClient {
	return g.pizza
}

func (g *GRPCClient) Cart() pb.CartServiceClient {
	return g.cart
}

func (g *GRPCClient) Order() pb.OrderServiceClient {
	return g.order
}

func NewGRPCClient() (ServiceManager, error) {
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		auth:  pb.NewAuthServiceClient(conn),
		pizza: pb.NewPizzaServiceClient(conn),
		cart:  pb.NewCartServiceClient(conn),
		order: pb.NewOrderServiceClient(conn),
	}, err
}
