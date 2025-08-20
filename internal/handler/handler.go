package handler

import (
	"github/http/copy/task4/config"
	grpc_client "github/http/copy/task4/internal/transport/grpc/client"
)

type Handler struct {
	GRPCClient grpc_client.ServiceManager
	cfg        config.Config
}

func NewHandler(GRPCCClient grpc_client.ServiceManager, cfg config.Config) *Handler {
	return &Handler{
		GRPCClient: GRPCCClient,
		cfg:       cfg,
	}
}
