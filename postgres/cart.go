package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/models"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type cart struct {
	db *sql.DB
}

type CartRepo interface {
	Cart(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error)
	CartItems(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error)
	GetCartId(ctx context.Context, userId int32) (*pizza.CartItemsResp, error)
	GetCartItemId(ctx context.Context, pizzaId int32, cartId int32) (*pizza.CartItemsResp, error)
	GetFromCart(ctx context.Context, id int32) (*pizza.CartItemsResp, error)
	CheckIsCartExist(ctx context.Context, req *pizza.CheckIsCartExistRequest) (*pizza.CheckIsCartExistResponse, error)
	IncreaseAmountOfPizza(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
	IncreaseTotalCost(ctx context.Context, id int32) (*pizza.CartItemsResp, error)
	DecreaseAmountOfPizza(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
	GetCartrHistory(ctx context.Context, req *pizza.GetCartHistoryRequest) (*pizza.GetCartHistoryResponse, error)
	GetCartItemHistory(ctx context.Context, req *pizza.GetCarItemtHistoryRequest) (*pizza.GetCarItemtHistoryResponse, error)
}

func NewCart(db *sql.DB) CartRepo {
	return &cart{
		db: db,
	}
}

func (c *cart) CheckIsCartExist(ctx context.Context, req *pizza.CheckIsCartExistRequest) (*pizza.CheckIsCartExistResponse, error) {
	var cart models.Cart

	query := `SELECT is_active FROM cart WHERE id = $1 AND user_id = $2`

	err := c.db.QueryRow(query, req.Id, req.UserId).Scan(&cart.IsActive)
	if err != nil {
		return nil, err
	}

	return &pizza.CheckIsCartExistResponse{
		Message:  "success",
		IsActive: cart.IsActive,
	}, nil
}

func (c *cart) Cart(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error) {

	query := `INSERT INTO cart(user_id, is_active, total_cost) 
	VALUES($1, $2, $3)`

	_, err := c.db.Exec(
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

func (c *cart) CartItems(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error) {

	query := `INSERT INTO cart_item(pizza_id, pizza_type_id, cost, cart_id, quantity)
	VALUES($1, $2, $3, $4, $5)`

	_, err := c.db.Exec(
		query,
		req.PizzaId,
		req.PizzaTypeId,
		req.Cost,
		req.Id,
		req.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.CartResponse{
		Message: "success",
	}, nil
}

func (c *cart) GetCartId(ctx context.Context, userId int32) (*pizza.CartItemsResp, error) {

	query := `
    SELECT id FROM cart WHERE user_id = $1
`
	var cartId int32

	err := c.db.QueryRow(query, userId).Scan(&cartId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &pizza.CartItemsResp{
		CartId: cartId,
	}, nil
}

func (c *cart) GetCartItemId(ctx context.Context, pizzaId int32, cartId int32) (*pizza.CartItemsResp, error) {

	query := `
    SELECT id FROM cart_item WHERE pizza_id = $1 AND cart_id = $2
`
	var id int32

	err := c.db.QueryRow(query, pizzaId, cartId).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &pizza.CartItemsResp{
		Id: id,
	}, nil
}

func (c *cart) GetFromCart(ctx context.Context, id int32) (*pizza.CartItemsResp, error) {

	query := `
    SELECT pizza_id, cart_id, quantity FROM cart_item WHERE id = $1
`
	var pizzaId, cartId, quantity int32
	err := c.db.QueryRow(query, id).Scan(&pizzaId, &cartId, &quantity)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &pizza.CartItemsResp{
		CartId:   cartId,
		PizzaId:  pizzaId,
		Quantity: quantity,
	}, nil
}

func (c *cart) IncreaseAmountOfPizza(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	query := `UPDATE cart_item SET quantity = $3, cost = $4 WHERE id = $1 AND pizza_id = $2`

	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		query,
		req.Id,
		req.PizzaId,
		req.Quantity,
		req.Cost,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &pizza.CartItemsResp{
		Message: "success",
	}, nil
}

func (c *cart) IncreaseTotalCost(ctx context.Context, id int32) (*pizza.CartItemsResp, error) {
	query := `
	UPDATE cart
	SET total_cost = (
    	SELECT COALESCE(SUM(ci.cost), 0)
    	FROM cart_item ci
 		WHERE ci.cart_id = $1
	)
	WHERE id = $1;
	`
	_, err := c.db.Exec(
		query,
		id,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.CartItemsResp{
		Message: "success",
	}, nil
}

func (c *cart) DecreaseAmountOfPizza(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	var query string

	if req.Quantity == 0 {
		query = `
	DELETE FROM cart_item WHERE id = $1
	`
	} else if req.Quantity > 1 {
		query = `
	UPDATE cart_item SET quantity = $1 WHERE id = $2
	`
	} else {
		return nil, errors.New("this pizza is bot exists in your cart")
	}
	_, err := c.db.Exec(
		query,
		req.Id,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.CartItemsResp{
		Message: "success",
	}, nil
}

func (c *cart) GetCartrHistory(ctx context.Context, req *pizza.GetCartHistoryRequest) (*pizza.GetCartHistoryResponse, error) {

	query := `SELECT c.id, c.is_active, o.date
          FROM cart AS c
          INNER JOIN orders AS o ON o.cart_id = c.id
          WHERE c.user_id = $1`

	rows, err := c.db.Query(query, req.UserId)
	if err != nil {
		return nil, err
	}

	var cartHistory []*pizza.GetCartHistoryResponse

	defer rows.Close()

	for rows.Next() {
		var cartId int32
		var isActive bool
		var date time.Time
		if err := rows.Scan(&cartId, &isActive, &date); err != nil {
			return nil, err
		}

		cartHistory = append(cartHistory, &pizza.GetCartHistoryResponse{
			CartId:   cartId,
			IsActive: isActive,
			Date:     timestamppb.New(date),
		})
	}

	return &pizza.GetCartHistoryResponse{
		CartHistory: cartHistory,
	}, nil
}

func (c *cart) GetCartItemHistory(ctx context.Context, req *pizza.GetCarItemtHistoryRequest) (*pizza.GetCarItemtHistoryResponse, error) {

	var cart models.CartIeamHistory

	query := `SELECT ci.pizza_id, ci.pizza_type_id, ci.cost, ci.quantity, c.total_cost
	FROM cart_item AS ci 
	JOIN cart AS c ON c.id = ci.cart_id
	WHERE c.id = $1`

	err := c.db.QueryRow(query, req.CartHistoryId).Scan(&cart.PizzaId, &cart.PizzaTypeId, &cart.Cost, &cart.Quantity, &cart.TotalCost)
	if err != nil {
		return nil, err
	}

	return &pizza.GetCarItemtHistoryResponse{
		CartHistoryId: req.CartHistoryId,
		PizzaId:       cart.PizzaId,
		PizzaTypeId:   cart.PizzaTypeId,
		Cost:          cart.Cost,
		Quantity:      cart.Quantity,
	}, nil
}
