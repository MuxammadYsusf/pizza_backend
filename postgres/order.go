package postgres

import (
	"context"
	"database/sql"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/models"
)

type orders struct {
	db *sql.DB
}
type OrderRepo interface {
	Order(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	OrderItem(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	CheckIsOrdered(ctx context.Context, req *pizza.CheckIsOrderedRequest) (*pizza.CheckIsOrderedResponse, error)
}

func NewOrder(db *sql.DB) OrderRepo {
	return &orders{
		db: db,
	}
}

func (o *orders) CheckIsOrdered(ctx context.Context, req *pizza.CheckIsOrderedRequest) (*pizza.CheckIsOrderedResponse, error) {
	var cart models.Order

	query := `SELECT is_ordered, status FROM orders WHERE id = $1 AND user_id = $2`

	err := o.db.QueryRow(query, req.Id, req.UserId).Scan(&cart.IsOrdered, &cart.Status)
	if err != nil {
		return nil, err
	}

	return &pizza.CheckIsOrderedResponse{
		Message:   "success",
		IsOrdered: cart.IsOrdered,
	}, nil
}

func (o *orders) Order(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	query := `INSERT INTO orders(date, is_ordered, user_id, status) 
	VALUES($1, $2, $3, $4)`

	_, err := o.db.Exec(
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

func (o *orders) OrderItem(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	query := `INSERT INTO order_item(pizza_id, pizza_type_id, cart_id, quantity) 
	VALUES($1, $2, $3, $4)`

	_, err := o.db.Exec(
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
