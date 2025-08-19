package handler

import (
	"database/sql"
	"github/http/copy/task4/grpc/client"
	"github/http/copy/task4/postgres"
)

type Handler struct {
	GRPCClient client.ServiceManager
	DB         postgres.NewPostgresI
}

func NewHandler(GRPCCLient client.ServiceManager, db *sql.DB) *Handler {
	return &Handler{
		GRPCClient: GRPCCLient,
		DB:         postgres.NewPostgres(db),
	}
}
