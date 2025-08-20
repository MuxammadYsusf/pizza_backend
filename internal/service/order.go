package service

import (
	"context"
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/config"
	"github/http/copy/task4/internal/storage"
	client "github/http/copy/task4/internal/transport/grpc/client"
	"github/http/copy/task4/pkg/util"
	"time"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type orderService struct {
	cfg      config.Config
	log      logger.LoggerI
	storage  storage.Storage
	services client.ServiceManager
	pb.UnimplementedOrderServiceServer
}

func NewOrderService(cfg config.Config, log logger.LoggerI, storage storage.Storage, services client.ServiceManager) *orderService {
	return &orderService{
		cfg:      cfg,
		log:      log,
		storage:  storage,
		services: services,
	}
}

func (s *orderService) OrderPizza(ctx context.Context, req *pb.OrderPizzaRequest) (*pb.OrderPizzaResponse, error) {

	var resp *pb.OrderPizzaResponse

	exists, err := s.storage.Postgres().Order().CheckIsOrdered(ctx, &pb.CheckIsOrderedRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	if exists.IsOrdered || exists.Status == util.STATUS_IN_PROGRESS {
		req.Date = timestamppb.New(time.Now())
		resp = &pb.OrderPizzaResponse{
			Message: "already ordered",
			Status:  exists.Status,
			Date:    req.Date,
		}
	} else {

		if exists.Status == util.STATUS_DONE || exists.Status == util.STATUS_IN_PROGRESS {
			req.IsOrdered = true
			req.Status = exists.Status
		}

		cId, err := s.storage.Postgres().Cart().GetCartId(ctx, req.UserId)
		if err != nil {
			return nil, err
		}

		req.CartId = cId.CartId
		req.Date = timestamppb.New(time.Now())
		req.Status = util.STATUS_IN_PROGRESS
		req.IsOrdered = true

		_, err = s.storage.Postgres().Order().Order(ctx, req)
		if err != nil {
			return nil, err
		}

		resp, err := s.storage.Postgres().Order().GetOrderId(ctx, req)
		if err != nil {
			return nil, err
		}

		req.Id = resp.Id

		pc, err := s.storage.Postgres().Pizza().GetAllPizzaCost(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		req.Cost = pc.Cost

		resp, err = s.storage.Postgres().Order().GetOrderItemId(ctx, req)
		if err != nil {
			return nil, err
		}

		resp, err = s.storage.Postgres().Order().OrderItem(ctx, req)
		if err != nil {
			return nil, err
		}

		_, err = s.storage.Postgres().Order().UpdateOrderStatus(ctx, req)
		if err != nil {
			return nil, err
		}

		isActive := false

		_, err = s.storage.Postgres().Cart().CloseTheCart(ctx, req.CartId, isActive)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
