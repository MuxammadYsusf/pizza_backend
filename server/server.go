package server

import (
	"database/sql"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/generated/session"

	client "github/http/copy/task4/grpc/client"
	service "github/http/copy/task4/grpc/service"

	"google.golang.org/grpc"
)

func NewServer(db *sql.DB, serviceManager client.ServiceManager) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer()

	session.RegisterAuthServiceServer(grpcServer, service.NewLoginService(db, serviceManager))
	pizza.RegisterPizzaServiceServer(grpcServer, service.NewPizzaService(db, serviceManager))

	return

}
