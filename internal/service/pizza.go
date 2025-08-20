package service

import (
	"context"
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/config"
	"github/http/copy/task4/internal/storage"
	client "github/http/copy/task4/internal/transport/grpc/client"

	"github.com/saidamir98/udevs_pkg/logger"
)

type pizzaService struct {
	cfg      config.Config
	log      logger.LoggerI
	storage  storage.Storage
	services client.ServiceManager
	pb.UnimplementedPizzaServiceServer
}

func NewPizzaService(cfg config.Config, log logger.LoggerI, storage storage.Storage, services client.ServiceManager) *pizzaService {
	return &pizzaService{
		cfg:      cfg,
		log:      log,
		storage:  storage,
		services: services,
	}
}

func (s *pizzaService) CreatePizzaType(ctx context.Context, req *pb.CreatePizzaRequest) (*pb.CreatePizzaResponse, error) {

	resp, err := s.storage.Postgres().Pizza().CreatePizzaType(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *pizzaService) CreatePizza(ctx context.Context, req *pb.CreatePizzaRequest) (*pb.CreatePizzaResponse, error) {

	resp, err := s.storage.Postgres().Pizza().CreatePizza(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
func (s *pizzaService) GetPizzas(ctx context.Context, req *pb.GetPizzasRequest) (*pb.GetPizzasResponse, error) {

	resp, err := s.storage.Postgres().Pizza().GetPizzas(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *pizzaService) GetPizzaById(ctx context.Context, req *pb.GetPizzaByIdRequest) (*pb.GetPizzaByIdResponse, error) {

	resp, err := s.storage.Postgres().Pizza().GetPizzaById(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *pizzaService) UpdatePizza(ctx context.Context, req *pb.UpdatePizzaRequest) (*pb.UpdatePizzaResponse, error) {

	_, err := s.storage.Postgres().Pizza().GetPizzaData(ctx, &pb.CartItems{
		PizzaId: req.Id,
	})
	if err != nil {
		return nil, err
	}

	resp, err := s.storage.Postgres().Pizza().UpdatePizza(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *pizzaService) DeletePizza(ctx context.Context, req *pb.DeletePizzaRequest) (*pb.DeletePizzaResponse, error) {

	resp, err := s.storage.Postgres().Pizza().DeletePizza(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *pizzaService) GetPizzaCost(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error) {

	resp, err := s.storage.Postgres().Pizza().GetPizzaCost(ctx, req.PizzaId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
