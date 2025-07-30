package service

import (
	"context"
	"github/http/copy/task4/generated/pizza"
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
