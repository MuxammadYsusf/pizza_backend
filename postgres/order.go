package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	c "github/http/copy/task4/constants"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/models"
	"strings"

	"github.com/lib/pq"
)

type orders struct {
	db *sql.DB
}
type OrderRepo interface {
	Order(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	OrderItem(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	CheckIsOrdered(ctx context.Context, req *pizza.CheckIsOrderedRequest) (*pizza.CheckIsOrderedResponse, error)
	GetOrderId(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	GetOrderItemId(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	UpdateOrderStatus(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaRequest, error)
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
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if cart.Status == c.OrderStatusDone || cart.Status == c.OrderStatusInProgress {
		cart.IsOrdered = true
	}

	return &pizza.CheckIsOrderedResponse{
		Message:   "success",
		IsOrdered: cart.IsOrdered,
		Status:    cart.Status,
	}, nil
}

func (o *orders) Order(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	query := `SELECT cart_id FROM orders WHERE id = $1
	`

	var cartId int32

	err := o.db.QueryRow(query, req.CartId).Scan(&cartId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if req.CartId == cartId {
		return &pizza.OrderPizzaResponse{
			Message: "already ordered",
		}, nil
	} else {
		query = `INSERT INTO orders(date, is_ordered, user_id, status, cart_id) 
	VALUES($1, $2, $3, $4, $5)`

		timeVal := req.Date.AsTime()

		_, err = o.db.Exec(
			query,
			timeVal,
			req.IsOrdered,
			req.UserId,
			req.Status,
			req.CartId,
		)
		if err != nil {
			return nil, err
		}
	}

	return &pizza.OrderPizzaResponse{
		Message: "success",
	}, nil
}

func (o *orders) GetOrderId(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {
	var order models.Order

	query := `SELECT id FROM orders WHERE id = $1 AND user_id = $2 ORDER BY id DESC LIMIT 1`

	err := o.db.QueryRow(query, req.CartId, req.UserId).Scan(&order.ID)
	if err != nil {
		return nil, err
	}

	fmt.Println("order", order.ID)

	return &pizza.OrderPizzaResponse{
		Message: "success",
		Id:      order.ID,
	}, nil
}

// DROP
func (o *orders) GetOrderItemId(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	query := `SELECT array_agg(id ORDER BY id DESC)
	FROM (
    	SELECT id
    	FROM order_item
    	WHERE order_id = $1
    	ORDER BY id DESC
    	LIMIT $2
	) sub;`

	var ids []int32
	var exists bool

	err := o.db.QueryRow(query, req.Id, req.Limit).Scan(pq.Array(&ids))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if len(ids) != 0 {
		exists = true
		fmt.Println("exists", exists)
	}

	return &pizza.OrderPizzaResponse{
		Message: "success",
		ItemIds: ids,
		IsExist: exists,
	}, nil
}

func (o *orders) OrderItem(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	tx, err := o.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Commit()

	query := `SELECT array_agg(id ORDER BY id DESC)
	FROM (
    	SELECT id
    	FROM order_item
    	WHERE order_id = $1
    	ORDER BY id DESC
    	LIMIT $2
	) sub;`

	var ids []int32
	var exists bool

	err = o.db.QueryRow(query, req.Id, req.Limit).Scan(pq.Array(&ids))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if len(ids) != 0 {
		exists = true
		fmt.Println("exists", exists)
	}

	if exists == true {
		return &pizza.OrderPizzaResponse{
			Message: "already ordered",
		}, nil
	}

	query = `INSERT INTO order_item(pizza_id, cost, quantity, order_id) VALUES`
	values := []interface{}{}
	placeholders := []string{}

	for i, items := range req.Items {

		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d, $%d)",
				i*4+1, i*4+2, i*4+3, i*4+4))
		values = append(values,
			items.PizzaId,
			req.Cost[i],
			items.Quantity,
			req.Id,
		)
	}

	if len(placeholders) == 0 {
		return nil, errors.New("no items")
	}

	query += strings.Join(placeholders, ",") + ";"

	_, err = tx.Exec(query, values...)
	if err != nil {
		fmt.Println("err", err)
		tx.Rollback()
		return nil, err
	}

	return &pizza.OrderPizzaResponse{
		Message: "success",
	}, nil
}

func (o *orders) UpdateOrderStatus(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaRequest, error) {

	query := `UPDATE orders SET status = $1, is_ordered = $2 WHERE id = $3`

	tx, err := o.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		query,
		req.Status,
		req.IsOrdered,
		req.Id,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &pizza.OrderPizzaRequest{}, nil
}
