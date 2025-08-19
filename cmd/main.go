package main

import (
	"github/http/copy/task4/config"
	"github/http/copy/task4/grpc/client"
	"github/http/copy/task4/handler"
	"github/http/copy/task4/pkg/db"
	"github/http/copy/task4/server"

	"net"

	"github.com/gin-contrib/cors"
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

	cont := handler.NewHandler(clients, conn)

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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5179"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/pizzas/register", cont.Register)
	r.POST("/pizzas/login", cont.Login)

	auth := r.Group("/", cont.AuthMiddleware)
	{
		auth.GET("/pizzas/get", cont.GetPizzas)
		auth.GET("/pizzas/get/:id/:typeId", cont.GetPizzaById)
		auth.POST("/pizzas/cart", cont.PutPizzaIntoCart)
		auth.PUT("/pizzas/decrease", cont.DecreasePizzaQuantity)
		auth.GET("/pizzas/get/pizza-cost/:id", cont.GetPizzaCost)
		auth.GET("/pizzas/get/total-cost/:id", cont.GetTotalCost)
		auth.GET("/pizzas/cart", cont.GetFromCart)
		auth.DELETE("/pizzas/cart/:pizzaId", cont.ClearTheCartById)
		auth.DELETE("/pizzas/cart", cont.ClearTheCart)
		auth.POST("pizzas/order", cont.OrderPizza)
		auth.GET("/pizzas/history", cont.GetCartHistory)
		auth.GET("/pizzas/history/:id", cont.GetCartItemHistory)
		auth.PUT("/pizzas/logout", cont.Logout)
	}

	admin := r.Group("/admin", cont.AuthMiddleware)
	{
		admin.POST("/pizzas/create/type", cont.CreatePizzaType)
		admin.POST("/pizzas/create", cont.CreatePizza)
		admin.GET("/pizzas/get", cont.GetPizzas)
		admin.GET("/pizzas/get/:id/:typeId", cont.GetPizzaById)
		admin.PUT("/pizzas/update", cont.UpdatePizza)
		admin.DELETE("/pizzas/delete/:id", cont.DeletePizza)
		admin.POST("/pizzas/cart", cont.PutPizzaIntoCart)
		admin.PUT("/pizzas/decrease", cont.DecreasePizzaQuantity)
		admin.GET("/pizzas/get/pizza-cost/:id", cont.GetPizzaCost)
		admin.GET("/pizzas/get/total-cost/:id", cont.GetTotalCost)
		admin.GET("/pizzas/cart", cont.GetFromCart)
		admin.DELETE("/pizzas/cart/:pizzaId", cont.ClearTheCartById)
		admin.DELETE("/pizzas/cart", cont.ClearTheCart)
		admin.POST("pizzas/order", cont.OrderPizza)
		admin.GET("/pizzas/history", cont.GetCartHistory)
		admin.GET("/pizzas/history/:id", cont.GetCartItemHistory)
		admin.PUT("/pizzas/logout", cont.Logout)
	}

	r.Run(cfg.HttpPort)
}
