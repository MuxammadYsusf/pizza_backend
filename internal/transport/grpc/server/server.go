package grpc_server

import (
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/config"
	"github/http/copy/task4/internal/service"
	"github/http/copy/task4/internal/storage"

	client "github/http/copy/task4/internal/transport/grpc/client"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServerParams struct {
	Cfg      *config.Config
	Store    storage.Storage
	Services client.ServiceManager
	Log      logger.LoggerI
}

func SetUpServer(cfg *config.Config, log logger.LoggerI, strg storage.Storage, svcs client.ServiceManager) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer()

	reflection.Register(grpcServer)

	return
}

func New(params GrpcServerParams) *grpc.Server {
	// authMiddleware := middleware.NewAuthMiddleware(params.Clients.AuthServiceClient)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
			// authMiddleware.Auth,
			// middleware.GrpcLoggerMiddleware,
			// middleware.GrpcErrorMiddleware,
			),
		),
	)

	reflection.Register(grpcServer)

	pb.RegisterAuthServiceServer(grpcServer, service.NewAuthService(*params.Cfg, params.Log, params.Store, params.Services))
	pb.RegisterPizzaServiceServer(grpcServer, service.NewPizzaService(*params.Cfg, params.Log, params.Store, params.Services))
	pb.RegisterCartServiceServer(grpcServer, service.NewCartService(*params.Cfg, params.Log, params.Store, params.Services))
	pb.RegisterOrderServiceServer(grpcServer, service.NewOrderService(*params.Cfg, params.Log, params.Store, params.Services))

	return grpcServer
}
