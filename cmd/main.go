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

	r.POST("/pizzas/register", cont.Register)
	r.POST("/pizzas/login", cont.Login)

	auth := r.Group("/", cont.AuthMiddleware)
	{
		auth.GET("/pizzas/get", cont.GetPizzas)
		auth.GET("/pizzas/get/:id/:typeId", cont.GetPizzaById)
		auth.POST("/pizzas/cart", cont.PutPizzaIntoCart)
		auth.POST("pizzas/order", cont.OrderPizza)
		auth.PUT("pizzas/addpizza:id/:pizzaId", cont.IncreasePizzaInCart)
		auth.PUT("/pizzas/removepizza", cont.DecreasePizzaInCart)
		auth.GET("/pizzas/history", cont.GetCartHistory)
		auth.GET("/pizzas/history/:id", cont.GetCartItemHistory)
	}

	admin := r.Group("/admin", cont.AdminOnlyMiddleware)
	{
		admin.POST("/pizzas/create/type", cont.CreatePizzaType)
		admin.POST("/pizzas/create", cont.CreatePizza)
		admin.GET("/pizzas/get", cont.GetPizzas)
		admin.GET("/pizzas/get/:id/:typeId", cont.GetPizzaById)
		admin.PUT("/pizzas/update/", cont.UpdatePizza)
		admin.DELETE("/pizzas/delete/:id/:typeId", cont.DeletePizza)
		admin.POST("/pizzas/cart/:id", cont.PutPizzaIntoCart)
		admin.PUT("pizzas/addpizza/:id/:pizzaId", cont.IncreasePizzaInCart)
		admin.PUT("/pizzas/removepizza", cont.DecreasePizzaInCart)
		admin.GET("/pizzas/history", cont.GetCartHistory)
		admin.GET("/pizzas/history/:id", cont.GetCartItemHistory)
	}

	r.Run(cfg.HttpPort)
}
