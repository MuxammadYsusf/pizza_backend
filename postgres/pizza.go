package postgres

import (
	"context"
	"database/sql"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/models"
)

type pizzas struct {
	db *sql.DB
}

type PizzaRepo interface {
	CreatePizza(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error)
	GetPizzaById(ctx context.Context, req *pizza.GetPizzaByIdRequest) (*pizza.GetPizzaByIdResponse, error)
	GetPizzas(ctx context.Context, req *pizza.GetPizzasRequest) (*pizza.GetPizzasResponse, error)
	UpdatePizza(ctx context.Context, req *pizza.UpdatePizzaRequest) (*pizza.UpdatePizzaResponse, error)
	DeletePizza(ctx context.Context, req *pizza.DeletePizzaRequest) (*pizza.DeletePizzaResponse, error)
	CheckIsCartExist(ctx context.Context, req *pizza.CheckIsCartExistRequest) (*pizza.CheckIsCartExistResponse, error)
	Cart(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error)
	CartItems(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error)
	CheckIsOrdered(ctx context.Context, req *pizza.CheckIsOrderedRequest) (*pizza.CheckIsOrderedResponse, error)
	Order(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	OrderItem(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
}

func NewPizza(db *sql.DB) PizzaRepo {
	return &pizzas{
		db: db,
	}
}

func (p *pizzas) CreatePizza(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error) {

	query := `INSERT INTO pizza VALUES(name, price, id, type_id)`

	_, err := p.db.Exec(
		query,
		req.Name,
		req.Price,
		req.Id,
		req.TypeId,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.CreatePizzaResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) GetPizzaById(ctx context.Context, req *pizza.GetPizzaByIdRequest) (*pizza.GetPizzaByIdResponse, error) {
	var pizzas models.Pizza

	query := `SELECT name, price FROM pizza WHERE id = $1 AND type_id= $2`

	err := p.db.QueryRow(query, req.Id, req.TypeId).Scan(req.Id, req.TypeId, pizzas.Name, pizzas.Price)
	if err != nil {
		return nil, err
	}

	return &pizza.GetPizzaByIdResponse{
		Name:  pizzas.Name,
		Price: pizzas.Price,
	}, nil
}

func (p *pizzas) GetPizzas(ctx context.Context, req *pizza.GetPizzasRequest) (*pizza.GetPizzasResponse, error) {

	query := `SELECT name, price FROM pizza`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pizzas []*pizza.Pizzas

	for rows.Next() {
		var name string
		var price float32
		if err := rows.Scan(name, price); err != nil {
			return nil, err
		}

		pizzas = append(pizzas, &pizza.Pizzas{
			Name:  name,
			Price: price,
		})
	}

	return &pizza.GetPizzasResponse{
		Pizzas: pizzas,
	}, nil
}

func (p *pizzas) UpdatePizza(ctx context.Context, req *pizza.UpdatePizzaRequest) (*pizza.UpdatePizzaResponse, error) {

	query := `UPDATE pizza SET name = $3, price = $4 WHERE id = $1 AND type_id = $2`

	_, err := p.db.Exec(
		query,
		req.Name,
		req.Price,
		req.Id,
		req.TypeId,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.UpdatePizzaResponse{
		Message: "success",
		Name:    req.Name,
		Price:   req.Price,
	}, nil
}

func (p *pizzas) DeletePizza(ctx context.Context, req *pizza.DeletePizzaRequest) (*pizza.DeletePizzaResponse, error) {

	query := `DELETE FROM pizza WHERE id = $1 AND typeId = $2`

	result, err := p.db.Exec(
		query,
		req.Id,
		req.TypeId,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	return &pizza.DeletePizzaResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) CheckIsCartExist(ctx context.Context, req *pizza.CheckIsCartExistRequest) (*pizza.CheckIsCartExistResponse, error) {
	var cart models.Cart

	query := `SELECT is_acive, total_cost FROM cart WHERE id = $1 AND user_id = $2`

	err := p.db.QueryRow(query, req.Id, req.UserId).Scan(req.Id, cart.IsActive, cart.TotalCost)
	if err != nil {
		return nil, err
	}

	return &pizza.CheckIsCartExistResponse{
		Message:  "success",
		IsActive: cart.IsActive,
	}, nil
}

func (p *pizzas) Cart(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error) {

	query := `INSERT INTO cart VALUES(user_id, is_active, tootal_cost)`

	_, err := p.db.Exec(
		query,
		req.UserId,
		req.IsActive,
		req.TotalCost,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.CartResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) CartItems(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error) {

	query := `INSERT INTO cart_item VALUES(pizza_id, pizza_type_id, cost, cart_id, quantity)`

	_, err := p.db.Exec(
		query,
		req.PizzaId,
		req.PizzaTypeId,
		req.Cost,
		req.CartId,
		req.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.CartResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) CheckIsOrdered(ctx context.Context, req *pizza.CheckIsOrderedRequest) (*pizza.CheckIsOrderedResponse, error) {
	var cart models.Order

	query := `SELECT is_ordered, date FROM orders WHERE id = $1 AND user_id = $2`

	err := p.db.QueryRow(query, req.Id, req.UserId).Scan(req.Id, cart.IsOrdered, cart.Date)
	if err != nil {
		return nil, err
	}

	return &pizza.CheckIsOrderedResponse{
		Message:   "success",
		IsOrdered: cart.IsOrdered,
	}, nil
}

func (p *pizzas) Order(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	query := `INSERT INTO orders VALUES(date, is_ordered, user_id, status)`

	_, err := p.db.Exec(
		query,
		req.Date,
		req.IsOrdered,
		req.UserId,
		req.Status,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.OrderPizzaResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) OrderItem(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	query := `INSERT INTO order_item VALUES(pizza_id, pizza_type_id, cart_id, quantity)`

	_, err := p.db.Exec(
		query,
		req.PizzaId,
		req.PizzaTypeId,
		req.CartId,
		req.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.OrderPizzaResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) GetUserHistory(ctx context.Context, req *pizza.GetUserHistoryRequest) (*pizza.GetUserHistoryResponse, error) {
	var cart models.Order

	query := `SELECT is_ordered, date FROM orders WHERE id = $1 AND user_id = $2`

	err := p.db.QueryRow(query, req.Id, req.UserId).Scan(req.Id, cart.IsOrdered, cart.Date)
	if err != nil {
		return nil, err
	}

	return &pizza.CheckIsOrderedResponse{
		Message:   "success",
		IsOrdered: cart.IsOrdered,
	}, nil
}
