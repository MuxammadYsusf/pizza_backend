package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/models"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type cart struct {
	db *sql.DB
}

type CartRepo interface {
	Cart(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error)
	CartItems(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error)
	IncreasePizzaQuantity(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error)
	DecreasePizzaQuantity(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error)
	GetCartId(ctx context.Context, userId int32) (*pb.CartItemsResp, error)
	GetCartItemId(ctx context.Context, pizzaId int32, cartId int32) (*pb.CartItemsResp, error)
	GetFromCart(ctx context.Context, Id int32, pizzaId int32) (*pb.CartItemsResp, error)
	CheckIsCartExist(ctx context.Context, req *pb.CheckIsCartExistRequest) (*pb.CheckIsCartExistResponse, error)
	ListCartItems(ctx context.Context, cartId int32) ([]*pb.CartItems, error)
	ClearTheCartById(ctx context.Context, cartId int32, pizzaId int32) (*pb.CartItemsResp, error)
	ClearTheCart(ctx context.Context, cartId int32) (*pb.CartItemsResp, error)
	GetTotalCost(ctx context.Context, id int32) (*pb.CartItemsResp, error)
	CloseTheCart(ctx context.Context, cartId int32, isActive bool) (*pb.CartResponse, error)
	GetCartHistory(ctx context.Context, req *pb.GetCartHistoryRequest) (*pb.GetCartHistoryResponse, error)
	GetCartItemHistory(ctx context.Context, req *pb.GetCartItemHistoryRequest) (*pb.GetCartItemHistoryResponse, error)
}

func NewCart(db *sql.DB) CartRepo {
	return &cart{
		db: db,
	}
}

func (c *cart) CheckIsCartExist(ctx context.Context, req *pb.CheckIsCartExistRequest) (*pb.CheckIsCartExistResponse, error) {
	var cart models.Cart

	query := `SELECT is_active FROM cart WHERE id = $1 AND user_id = $2`

	err := c.db.QueryRow(query, req.Id, req.UserId).Scan(&cart.IsActive)
	if err != nil {
		return nil, err
	}

	return &pb.CheckIsCartExistResponse{
		Message:  "success",
		IsActive: cart.IsActive,
	}, nil
}

func (c *cart) Cart(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error) {

	query := `INSERT INTO cart(user_id, is_active) 
	VALUES($1, $2)`

	_, err := c.db.Exec(
		query,
		req.UserId,
		req.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CartResponse{
		Message: "success",
	}, nil
}

func (c *cart) CartItems(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error) {

	query := `INSERT INTO cart_item(pizza_id, pizza_type_id, cart_id, quantity)
	VALUES($1, $2, $3, $4)`

	_, err := c.db.Exec(
		query,
		req.PizzaId,
		req.PizzaTypeId,
		req.Id,
		req.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CartResponse{
		Message: "success",
	}, nil
}

func (c *cart) IncreasePizzaQuantity(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error) {

	query := `UPDATE cart_item SET quantity = $3 WHERE id = $1 AND pizza_id = $2`

	_, err := c.db.Exec(
		query,
		req.CartItemId,
		req.PizzaId,
		req.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CartResponse{
		Message: "success",
	}, nil
}

func (c *cart) DecreasePizzaQuantity(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error) {

	var query string

	if req.Quantity == 0 {
		query = `
	DELETE FROM cart_item WHERE id = $1
	`
		_, err := c.db.Exec(
			query,
			req.Id,
		)
		if err != nil {
			return nil, err
		}

	} else if req.Quantity > 1 {
		query = `
	UPDATE cart_item SET quantity = $1 WHERE id = $2
	`
		_, err := c.db.Exec(
			query,
			req.Quantity,
			req.Id,
		)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, errors.New("this pizza is bot exists in your cart")
	}

	return &pb.CartItemsResp{
		Message: "success",
	}, nil
}

func (c *cart) GetCartId(ctx context.Context, userId int32) (*pb.CartItemsResp, error) {

	query := `
    SELECT id FROM cart WHERE user_id = $1 ORDER BY id DESC LIMIT 1
`
	var cartId int32

	err := c.db.QueryRow(query, userId).Scan(&cartId)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("HERE???", err)
		return nil, err
	}

	return &pb.CartItemsResp{
		CartId: cartId,
	}, nil
}

func (c *cart) GetCartItemId(ctx context.Context, pizzaId int32, cartId int32) (*pb.CartItemsResp, error) {

	query := `
    SELECT id FROM cart_item WHERE pizza_id = $1 AND cart_id = $2
`
	var id int32

	err := c.db.QueryRow(query, pizzaId, cartId).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &pb.CartItemsResp{
		Id: id,
	}, nil
}

func (c *cart) GetFromCart(ctx context.Context, Id int32, pizzaId int32) (*pb.CartItemsResp, error) {

	query := `
    SELECT pizza_id, cart_id, quantity, pizza_type_id FROM cart_item WHERE pizza_id = $1 AND cart_id = $2
`
	var scanPizzaId, cartId, quantity, pizzaTypeId int32
	err := c.db.QueryRow(query, pizzaId, Id).Scan(&scanPizzaId, &cartId, &quantity, &pizzaTypeId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &pb.CartItemsResp{
		CartId:      cartId,
		PizzaId:     scanPizzaId,
		Quantity:    quantity,
		PizzaTypeId: pizzaTypeId,
	}, nil
}

func (c *cart) ListCartItems(ctx context.Context, cartId int32) ([]*pb.CartItems, error) {

	const q = `
        SELECT
            ci.pizza_id,
            p.name,
            p.cost,
            p.photo,
            SUM(ci.quantity),
            MAX(ci.pizza_type_id) AS pizza_type_id
        FROM cart_item ci
        JOIN pizza p ON p.id = ci.pizza_id
        WHERE ci.cart_id = $1
        GROUP BY ci.pizza_id, p.name, p.cost, p.photo
        ORDER BY p.name;
    `
	rows, err := c.db.QueryContext(ctx, q, cartId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*pb.CartItems{}
	for rows.Next() {
		var it pb.CartItems
		if err := rows.Scan(&it.PizzaId, &it.Name, &it.Cost, &it.Photo, &it.Quantity, &it.PizzaTypeId); err != nil {
			return nil, err
		}
		items = append(items, &it)
	}
	return items, rows.Err()
}

func (c *cart) GetTotalCost(ctx context.Context, id int32) (*pb.CartItemsResp, error) {
	var total_cost float32

	query := `
	SELECT COALESCE(SUM(p.cost * ci.quantity), 0) AS total_cost
	FROM pizza p
	JOIN cart_item ci ON p.id = ci.pizza_id
	WHERE ci.cart_id = $1;
	`

	err := c.db.QueryRow(query, id).Scan(&total_cost)
	if err != nil {
		return nil, err
	}

	return &pb.CartItemsResp{
		TotalCost: total_cost,
	}, nil
}

func (c *cart) ClearTheCart(ctx context.Context, cartId int32) (*pb.CartItemsResp, error) {

	query := `DELETE FROM cart_item WHERE cart_id = $1`

	_, err := c.db.Exec(
		query,
		cartId,
	)

	if err != nil {
		return nil, err
	}

	return &pb.CartItemsResp{
		Message: "success",
	}, nil
}

func (c *cart) ClearTheCartById(ctx context.Context, cartId int32, pizzaId int32) (*pb.CartItemsResp, error) {

	query := `DELETE FROM cart_item WHERE cart_id = $1 AND pizza_id = $2`

	_, err := c.db.Exec(
		query,
		cartId,
		pizzaId,
	)

	if err != nil {
		return nil, err
	}

	return &pb.CartItemsResp{
		Message: "success",
	}, nil
}

func (c *cart) CloseTheCart(ctx context.Context, cartId int32, isActive bool) (*pb.CartResponse, error) {

	query := `UPDATE cart SET is_active = $1 WHERE id = $2`

	_, err := c.db.Exec(
		query,
		isActive,
		cartId,
	)

	if err != nil {
		return nil, err
	}

	return &pb.CartResponse{
		Message: "success",
	}, nil
}

func (c *cart) GetCartHistory(ctx context.Context, req *pb.GetCartHistoryRequest) (*pb.GetCartHistoryResponse, error) {

	query := `SELECT c.id, c.is_active, o.date
          FROM cart AS c
          INNER JOIN orders AS o ON o.cart_id = c.id
          WHERE c.user_id = $1`

	rows, err := c.db.Query(query, req.UserId)
	if err != nil {
		return nil, err
	}

	var cartHistory []*pb.GetCartHistoryResponse

	defer rows.Close()

	for rows.Next() {
		var cartId int32
		var isActive bool
		var date time.Time
		if err := rows.Scan(&cartId, &isActive, &date); err != nil {
			return nil, err
		}

		cartHistory = append(cartHistory, &pb.GetCartHistoryResponse{
			CartId:   cartId,
			IsActive: isActive,
			Date:     timestamppb.New(date),
		})
	}

	return &pb.GetCartHistoryResponse{
		CartHistory: cartHistory,
	}, nil
}

func (c *cart) GetCartItemHistory(ctx context.Context, req *pb.GetCartItemHistoryRequest) (*pb.GetCartItemHistoryResponse, error) {

	var cart models.CartItemHistory

	query := `SELECT oi.pizza_id, oi.cost, oi.quantity
	FROM order_item AS oi
	JOIN orders AS o ON o.id = oi.order_id
	WHERE o.id = $1`

	err := c.db.QueryRow(query, req.OrderId).Scan(&cart.PizzaId, &cart.Cost, &cart.Quantity)
	if err != nil {
		return nil, err
	}

	return &pb.GetCartItemHistoryResponse{
		CartHistoryId: req.OrderId,
		PizzaId:       cart.PizzaId,
		Cost:          cart.Cost,
		Quantity:      cart.Quantity,
	}, nil
}
