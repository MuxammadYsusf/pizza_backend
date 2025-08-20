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
		ois, err := s.pizzaPostgres.Order().CheckOrderItem(ctx, &pizza.CheckOrderItemRequest{
			OrderId: req.Id,
		})
		if err != nil {
			return nil, err
		}
		if len(ois.OrderItems) == 0 {
			resp, err = s.pizzaPostgres.Order().OrderItem(ctx, req)
			if err != nil {
				fmt.Println("err4", err)
				return nil, err
			}
		}

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
			fmt.Println("err1", err)
			return nil, err
		}

		req.Id = resp.Id

		pc, err := s.pizzaPostgres.Pizza().GetAllPizzaCost(ctx, req.CartId)
		if err != nil {
			fmt.Println("err2", err)
			return nil, err
		}

		req.Cost = pc.Cost

		resp, err = s.pizzaPostgres.Order().GetOrderItemId(ctx, req)
		if err != nil {
			fmt.Println("err3", err)
			return nil, err
		}

		resp, err = s.pizzaPostgres.Order().OrderItem(ctx, req)
		if err != nil {
			fmt.Println("err4", err)
			return nil, err
		}

		_, err = s.pizzaPostgres.Order().UpdateOrderStatus(ctx, req)
		if err != nil {
			return nil, err
		}

		isActive := false

		_, err = s.pizzaPostgres.Cart().CloseTheCart(ctx, req.CartId, isActive)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
