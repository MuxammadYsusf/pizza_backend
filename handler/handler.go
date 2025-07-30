package handler

import (
	"github/http/copy/task4/grpc/client"
)

type Handler struct {
	GRPCClient client.ServiceManager
}

func NewHandler(GRPCCLient client.ServiceManager) *Handler {
	return &Handler{
		GRPCClient: GRPCCLient,
	}
}
