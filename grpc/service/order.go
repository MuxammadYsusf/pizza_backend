package service

import (
	"context"
	"fmt"
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/pizza"
)

func (s *PizzaService) OrderPizza(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	exists, err := s.pizzaPostgres.Order().CheckIsOrdered(ctx, &pizza.CheckIsOrderedRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	if exists.IsOrdered || exists.Status == c.OrderStatusInProgress {
		return nil, fmt.Errorf("you already have an order in progress")
	}

	resp, err := s.pizzaPostgres.Order().Order(ctx, req)
	if err != nil {
		return nil, err
	}

	_, err = s.pizzaPostgres.Order().OrderItem(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
