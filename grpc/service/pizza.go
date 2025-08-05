package service

import (
	"context"
	"fmt"
	"github/http/copy/task4/generated/pizza"
)

const (
	statusInProgress = "in progress"
	statusDone       = "done"
	statusCanceled   = "canceled"
)

func (s *PizzaService) CreatePizza(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error) {

	resp, err := s.pizzaPostgres.Pizza().CreatePizza(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
func (s *PizzaService) GetPizzas(ctx context.Context, req *pizza.GetPizzasRequest) (*pizza.GetPizzasResponse, error) {

	resp, err := s.pizzaPostgres.Pizza().GetPizzas(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) GetPizzaById(ctx context.Context, req *pizza.GetPizzaByIdRequest) (*pizza.GetPizzaByIdResponse, error) {

	resp, err := s.pizzaPostgres.Pizza().GetPizzaById(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) UpdatePizza(ctx context.Context, req *pizza.UpdatePizzaRequest) (*pizza.UpdatePizzaResponse, error) {

	resp, err := s.pizzaPostgres.Pizza().UpdatePizza(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) DeletePizza(ctx context.Context, req *pizza.DeletePizzaRequest) (*pizza.DeletePizzaResponse, error) {

	resp, err := s.pizzaPostgres.Pizza().DeletePizza(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) Cart(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error) {

	var resp *pizza.CartResponse

	exists, err := s.pizzaPostgres.Pizza().CheckIsCartExist(ctx, &pizza.CheckIsCartExistRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	if !exists.IsActive {
		resp, err = s.pizzaPostgres.Pizza().Cart(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	resp = &pizza.CartResponse{
		Message: "the cart already exists",
	}

	_, err = s.pizzaPostgres.Pizza().CartItems(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) UpdatePizzaInCart(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	cart, err := s.pizzaPostgres.Pizza().GetFromCart(ctx, req)
	if err != nil {
		return nil, err
	}

	req.CartId = cart.CartId
	req.PizzaId = cart.PizzaId

	pizza, err := s.pizzaPostgres.Pizza().GetFromPizza(ctx, req)
	if err != nil {
		return nil, err
	}

	newCost := pizza.Cost * float32(req.Quantity)
	req.Cost = newCost

	resp, err := s.pizzaPostgres.Pizza().UpdatePizzaInCart(ctx, req)
	if err != nil {
		return nil, err
	}

	_, err = s.pizzaPostgres.Pizza().UpdateTotalCost(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) OrderPizza(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	exists, err := s.pizzaPostgres.Pizza().CheckIsOrdered(ctx, &pizza.CheckIsOrderedRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}
	if exists.IsOrdered || exists.Status == statusInProgress {
		return nil, fmt.Errorf("you already have an order in progress")
	}

	resp, err := s.pizzaPostgres.Pizza().Order(ctx, req)
	if err != nil {
		return nil, err
	}

	_, err = s.pizzaPostgres.Pizza().OrderItem(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) GetCartHistory(ctx context.Context, req *pizza.GetCartHistoryRequest) (*pizza.GetCartHistoryResponse, error) {

	resp, err := s.pizzaPostgres.Pizza().GetCartrHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) GetCartItemHistory(ctx context.Context, req *pizza.GetCarItemtHistoryRequest) (*pizza.GetCarItemtHistoryResponse, error) {

	resp, err := s.pizzaPostgres.Pizza().GetCartItemHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
