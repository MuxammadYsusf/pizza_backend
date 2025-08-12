package service

import (
	"context"
	"database/sql"
	"fmt"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/models"
)

func (s *PizzaService) Cart(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error) {

	var resp *pizza.CartResponse

	exists, err := s.pizzaPostgres.Cart().CheckIsCartExist(ctx, &pizza.CheckIsCartExistRequest{
		UserId: req.UserId,
		Id:     req.Id,
	})
	fmt.Println("cartId", req.Id)
	if err == sql.ErrNoRows || !exists.IsActive {
		req.IsActive = true
		req.TotalCost = 0
		resp, err = s.pizzaPostgres.Cart().Cart(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	resp = &pizza.CartResponse{
		Message: "the cart already exists",
	}

	var items models.Pizza

	items.ID = req.Items[0].PizzaId
	items.TypeId = req.Items[0].PizzaTypeId
	items.Quantity = req.Items[0].Quantity

	req.PizzaId = items.ID
	req.PizzaTypeId = items.TypeId
	req.Quantity = items.Quantity

	dataResp, err := s.pizzaPostgres.Pizza().GetPizzaData(ctx, &pizza.CartItems{
		UserId:  req.UserId,
		PizzaId: req.PizzaId,
	})
	if err != nil {
		return nil, err
	}

	cartItemId, err := s.pizzaPostgres.Cart().GetCartItemId(ctx, req.PizzaId, req.Id)
	if err != nil {
		return nil, err
	}

	req.CartItemId = cartItemId.Id

	cart, err := s.pizzaPostgres.Cart().GetFromCart(ctx, cartItemId.Id)
	if err != nil {
		return nil, err
	}

	newCost := dataResp.Cost * float32(req.Quantity)
	req.Cost = newCost

	if req.PizzaId == cart.PizzaId {

		req.Quantity = cart.Quantity + req.Quantity

		_, err := s.pizzaPostgres.Cart().IncreasePizzaQuantity(ctx, req)
		if err != nil {
			return nil, err
		}
	} else {

		_, err = s.pizzaPostgres.Cart().CartItems(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *PizzaService) DecreasePizzaQuantity(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	cartId, err := s.pizzaPostgres.Cart().GetCartId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	req.CartId = cartId.CartId

	cartItemId, err := s.pizzaPostgres.Cart().GetCartItemId(ctx, req.PizzaId, req.CartId)
	if err != nil {
		return nil, err
	}

	req.Id = cartItemId.Id

	q, err := s.pizzaPostgres.Cart().GetFromCart(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	req.Quantity = q.Quantity - req.Quantity

	resp, err := s.pizzaPostgres.Cart().DecreasePizzaQuantity(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) GetTotalCost(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	resp, err := s.pizzaPostgres.Cart().GetTotalCost(ctx, req.CartId)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) GetCartHistory(ctx context.Context, req *pizza.GetCartHistoryRequest) (*pizza.GetCartHistoryResponse, error) {

	resp, err := s.pizzaPostgres.Cart().GetCartrHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PizzaService) GetCartItemHistory(ctx context.Context, req *pizza.GetCarItemtHistoryRequest) (*pizza.GetCarItemtHistoryResponse, error) {

	resp, err := s.pizzaPostgres.Cart().GetCartItemHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
