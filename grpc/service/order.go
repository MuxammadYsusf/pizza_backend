package service

import (
	"context"
	"fmt"
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/pizza"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PizzaService) OrderPizza(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	var resp *pizza.OrderPizzaResponse

	exists, err := s.pizzaPostgres.Order().CheckIsOrdered(ctx, &pizza.CheckIsOrderedRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	if exists.IsOrdered || exists.Status == c.OrderStatusInProgress {
		req.Date = timestamppb.New(time.Now())
		resp = &pizza.OrderPizzaResponse{
			Message: "already ordered",
			Status:  exists.Status,
			Date:    req.Date,
		}
	} else {
		if exists.Status == c.OrderStatusDone || exists.Status == c.OrderStatusInProgress {
			req.IsOrdered = true
			req.Status = exists.Status
		}

		cId, err := s.pizzaPostgres.Cart().GetCartId(ctx, req.UserId)
		if err != nil {
			return nil, err
		}

		req.CartId = cId.CartId
		req.Date = timestamppb.New(time.Now())
		req.Status = c.OrderStatusInProgress
		req.IsOrdered = true

		_, err = s.pizzaPostgres.Order().Order(ctx, req)
		if err != nil {
			return nil, err
		}

		resp, err := s.pizzaPostgres.Order().GetOrderId(ctx, req)
		if err != nil {
			return nil, err
		}

		req.Id = resp.Id
		fmt.Println("Order Id", req.Id)

		tc, err := s.pizzaPostgres.Cart().GetTotalCost(ctx, req.UserId, req.CartId)
		if err != nil {
			return nil, err
		}

		req.TotalCost = tc.TotalCost
		req.Status = c.OrderStatusDone
		fmt.Println("Hi")

		resp, err = s.pizzaPostgres.Order().GetOrderItemId(ctx, req)
		if err != nil {
			return nil, err
		}

		req.ItemIds = resp.ItemIds

		_, err = s.pizzaPostgres.Order().OrderItem(ctx, req)
		if err != nil {
			return nil, err
		}

		_, err = s.pizzaPostgres.Order().UpdateOrderStatus(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
