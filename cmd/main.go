package main

import (
	"github/http/copy/task4/config"
	"github/http/copy/task4/grpc/client"
	"github/http/copy/task4/handler"
	"github/http/copy/task4/pkg/db"
	"github/http/copy/task4/server"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

func main() {
	cfg := config.Cfg()

	conn, err := db.Postgres(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	clients, err := client.NewGRPCClient()
	if err != nil {
		log.Fatal(err)
	}

	cont := handler.NewHandler(clients)

	grpcServer := server.NewServer(conn, clients)

	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatal(err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()

	r.POST("/login", cont.Login)
	r.POST("/register", cont.Register)

	auth := r.Group("/", cont.AuthMiddleware)
	{
		auth.GET("/pizzas", cont.GetPizzas)
		auth.GET("/pizzas/:id", cont.GetPizzaById)
		auth.POST("/cart", cont.PutPizzaIntoCart)
	}

	admin := r.Group("/admin", cont.AdminOnlyMiddleware)
	{
		admin.POST("/pizza", cont.CreatePizza)
		admin.GET("/pizzas", cont.GetPizzas)
		admin.GET("/pizzas/:id", cont.GetPizzaById)
		admin.PUT("/pizzas/:id", cont.UpdatePizza)
		admin.DELETE("/pizzas/:id", cont.DeletePizza)
		admin.POST("/cart", cont.PutPizzaIntoCart)
	}

	r.Run(cfg.HttpPort)
}
