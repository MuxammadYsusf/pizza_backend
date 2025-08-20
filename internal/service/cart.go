package service

import (
	"context"
	"database/sql"
	"fmt"
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/config"
	"github/http/copy/task4/internal/storage"
	client "github/http/copy/task4/internal/transport/grpc/client"
	"github/http/copy/task4/models"

	"github.com/saidamir98/udevs_pkg/logger"
)

type cartService struct {
	cfg      config.Config
	log      logger.LoggerI
	storage  storage.Storage
	services client.ServiceManager
	pb.UnimplementedCartServiceServer
}

func NewCartService(cfg config.Config, log logger.LoggerI, storage storage.Storage, services client.ServiceManager) *cartService {
	return &cartService{
		cfg:      cfg,
		log:      log,
		storage:  storage,
		services: services,
	}
}

func (s *cartService) Cart(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error) {

	var resp *pb.CartResponse

	ci, err := s.storage.Postgres().Cart().GetCartId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	req.Id = ci.CartId

	exists, err := s.storage.Postgres().Cart().CheckIsCartExist(ctx, &pb.CheckIsCartExistRequest{
		UserId: req.UserId,
		Id:     req.Id,
	})
	if err == sql.ErrNoRows || !exists.IsActive {
		req.IsActive = true
		req.TotalCost = 0
		resp, err = s.storage.Postgres().Cart().Cart(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	resp = &pb.CartResponse{
		Message: "the cart already exists",
	}

	var items models.Pizza

	items.ID = req.Items[0].PizzaId
	items.TypeId = req.Items[0].PizzaTypeId
	items.Quantity = req.Items[0].Quantity

	req.PizzaId = items.ID
	req.PizzaTypeId = items.TypeId
	req.Quantity = items.Quantity

	dataResp, err := s.storage.Postgres().Pizza().GetPizzaData(ctx, &pb.CartItems{
		UserId:  req.UserId,
		PizzaId: req.PizzaId,
	})
	if err != nil {
		return nil, err
	}

	cartItemId, err := s.storage.Postgres().Cart().GetCartItemId(ctx, req.PizzaId, req.Id)
	if err != nil {
		return nil, err
	}

	req.CartItemId = cartItemId.Id

	cart, err := s.storage.Postgres().Cart().GetFromCart(ctx, req.Id, req.PizzaId)
	if err != nil {
		return nil, err
	}

	newCost := dataResp.Cost * float32(req.Quantity)
	req.Cost = newCost

	if req.PizzaId == cart.PizzaId {

		req.Quantity = cart.Quantity + req.Quantity

		_, err := s.storage.Postgres().Cart().IncreasePizzaQuantity(ctx, req)
		if err != nil {
			return nil, err
		}
	} else {

		_, err = s.storage.Postgres().Cart().CartItems(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *cartService) DecreasePizzaQuantity(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error) {

	cartId, err := s.storage.Postgres().Cart().GetCartId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	req.CartId = cartId.CartId

	cartItemId, err := s.storage.Postgres().Cart().GetCartItemId(ctx, req.PizzaId, req.CartId)
	if err != nil {
		return nil, err
	}

	req.Id = cartItemId.Id

	q, err := s.storage.Postgres().Cart().GetFromCart(ctx, req.CartId, req.PizzaId)
	if err != nil {
		return nil, err
	}

	req.Quantity = q.Quantity - req.Quantity

	resp, err := s.storage.Postgres().Cart().DecreasePizzaQuantity(ctx, req)
	if err != nil {
		fmt.Println("err 3-->", err)
		return nil, err
	}

	return resp, nil
}

func (s *cartService) GetFromCart(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error) {

	cartId, err := s.storage.Postgres().Cart().GetCartId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	items, err := s.storage.Postgres().Cart().ListCartItems(ctx, cartId.CartId)
	if err != nil {
		return nil, err
	}

	var total float32
	for _, it := range items {
		total += it.Cost * float32(it.Quantity)
	}
	return &pb.CartItemsResp{CartItems: items, TotalCost: total}, nil
}

func (s *cartService) ClearTheCart(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error) {

	resp, err := s.storage.Postgres().Cart().GetCartId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	req.CartId = resp.CartId

	fmt.Println("req.CartId -->", req.CartId)

	resp, err = s.storage.Postgres().Cart().ClearTheCart(ctx, req.CartId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *cartService) ClearTheCartById(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error) {

	resp, err := s.storage.Postgres().Cart().GetCartId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	req.CartId = resp.CartId

	fmt.Println("req.CartId -->", req.CartId)

	resp, err = s.storage.Postgres().Cart().ClearTheCartById(ctx, req.CartId, req.PizzaId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *cartService) GetTotalCost(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error) {

	resp, err := s.storage.Postgres().Cart().GetTotalCost(ctx, req.CartId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *cartService) GetCartHistory(ctx context.Context, req *pb.GetCartHistoryRequest) (*pb.GetCartHistoryResponse, error) {

	resp, err := s.storage.Postgres().Cart().GetCartHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *cartService) GetCartItemHistory(ctx context.Context, req *pb.GetCartItemHistoryRequest) (*pb.GetCartItemHistoryResponse, error) {

	resp, err := s.storage.Postgres().Cart().GetCartItemHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
