package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github/http/copy/task4/config"
	"github/http/copy/task4/internal/handler"
	db "github/http/copy/task4/internal/storage"
	grpc_client "github/http/copy/task4/internal/transport/grpc/client"
	grpc_server "github/http/copy/task4/internal/transport/grpc/server"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/saidamir98/udevs_pkg/logger"
)

func main() {
	// Load environment variables from local files (non-fatal if absent).
	_ = godotenv.Load(".env", "env")

	// Load unified application configuration.
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize structured logger.
	log := logger.NewLogger(cfg.ServiceName, logger.LevelDebug)
	defer logger.Cleanup(log)

	// Initialize storage (PostgreSQL + sqlc). Panic if DB is not reachable.
	store, err := db.New(context.Background(), cfg, log)
	if err != nil {
		log.Panic("failed to connect to database", logger.Error(err))
	}
	defer store.Close()

	// Initialize outbound gRPC clients (if any external services are used).
	clients, err := grpc_client.NewGRPCClient()
	if err != nil {
		log.Panic("failed to create gRPC client", logger.Error(err))
	}

	// Build handler container (inject dependencies).
	cont := handler.NewHandler(clients, *cfg)

	// -------------------- HTTP (Gin) setup & route registration --------------------
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Restrict trusted proxies (avoid Gin warning; adjust if real proxies exist).
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Error("failed to set trusted proxies", logger.Error(err))
	}

	// CORS policy (adjust for production).
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5179"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes (no auth).
	r.POST("/pizzas/register", cont.Register)
	r.POST("/pizzas/login", cont.Login)

	// Authenticated user routes.
	auth := r.Group("", cont.AuthMiddleware)
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
		auth.POST("/pizzas/order", cont.OrderPizza)
		auth.GET("/pizzas/history", cont.GetCartHistory)
		auth.GET("/pizzas/history/:id", cont.GetCartItemHistory)
	}

	// Admin-only routes.
	admin := r.Group("/admin", cont.AuthMiddleware, cont.AdminOnlyMiddleware)
	{
		admin.POST("/pizzas/create/type", cont.CreatePizzaType)
		admin.POST("/pizzas/create", cont.CreatePizza)
		admin.GET("/pizzas/get", cont.GetPizzas)
		admin.GET("/pizzas/get/:id/:typeId", cont.GetPizzaById)
		admin.PUT("/pizzas/update", cont.UpdatePizza)
		admin.DELETE("/pizzas/delete/:id", cont.DeletePizza)
		admin.POST("/pizzas/order", cont.OrderPizza)
	}

	// HTTP server with sane timeouts.
	httpSrv := &http.Server{
		Addr:         cfg.HTTPPort, // Expected format ":8080"
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// -------------------- gRPC server (prepare, then start) --------------------
	grpcSrv := grpc_server.New(grpc_server.GrpcServerParams{
		Cfg:   cfg,
		Store: store,
		Log:   log,
	})

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Panic("failed to listen gRPC", logger.Error(err))
	}

	// -------------------- Ordered startup logging --------------------
	// 1. gRPC startup log (serve in goroutine).
	log.Info("gRPC listening", logger.String("port", cfg.GRPCPort))
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			// Normal graceful stop returns an internal gRPC error; ignore unless unexpected.
			if !errors.Is(err, context.Canceled) {
				log.Error("gRPC serve failed", logger.Error(err))
			}
		}
	}()

	// 2. HTTP startup log (serve in goroutine).
	log.Info("HTTP listening", logger.String("port", cfg.HTTPPort))
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error", logger.Error(err))
		}
	}()

	// -------------------- Signal handling (graceful shutdown) --------------------
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// First signal triggers graceful shutdown.
	sig := <-sigCh
	log.Info("received signal, shutting down servers...", logger.String("signal", sig.String()))

	// Second signal forces immediate exit.
	go func() {
		sig2 := <-sigCh
		log.Warn("second signal received, force exit", logger.String("signal", sig2.String()))
		os.Exit(1)
	}()

	// Gracefully stop gRPC server.
	grpcSrv.GracefulStop()

	// Gracefully stop HTTP server with timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP shutdown error", logger.Error(err))
		_ = httpSrv.Close()
	}

	log.Info("shutdown complete")
}
