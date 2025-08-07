package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/models"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type pizzas struct {
	db *sql.DB
}

type cart struct {
	db *sql.DB
}

type orders struct {
	db *sql.DB
}

type PizzaRepo interface {
	CreatePizzaType(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error)
	CreatePizza(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error)
	GetPizzaById(ctx context.Context, req *pizza.GetPizzaByIdRequest) (*pizza.GetPizzaByIdResponse, error)
	GetPizzas(ctx context.Context, req *pizza.GetPizzasRequest) (*pizza.GetPizzasResponse, error)
	UpdatePizza(ctx context.Context, req *pizza.UpdatePizzaRequest) (*pizza.UpdatePizzaResponse, error)
	DeletePizza(ctx context.Context, req *pizza.DeletePizzaRequest) (*pizza.DeletePizzaResponse, error)
	GetPizzaCost(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
}

type CartRepo interface {
	Cart(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error)
	CartItems(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error)
	GetFromCart(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
	CheckIsCartExist(ctx context.Context, req *pizza.CheckIsCartExistRequest) (*pizza.CheckIsCartExistResponse, error)
	IncreaseAmountOfPizza(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
	IncreaseTotalCost(ctx context.Context, id int32) (*pizza.CartItemsResp, error)
	DecreaseAmountOfPizza(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
	DecreaseTotalCost(ctx context.Context, id int32) (*pizza.CartItemsResp, error)
	GetCartrHistory(ctx context.Context, req *pizza.GetCartHistoryRequest) (*pizza.GetCartHistoryResponse, error)
	GetCartItemHistory(ctx context.Context, req *pizza.GetCarItemtHistoryRequest) (*pizza.GetCarItemtHistoryResponse, error)
}

type OrderRepo interface {
	Order(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	OrderItem(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	CheckIsOrdered(ctx context.Context, req *pizza.CheckIsOrderedRequest) (*pizza.CheckIsOrderedResponse, error)
}

func NewPizza(db *sql.DB) PizzaRepo {
	return &pizzas{
		db: db,
	}
}

func NewCart(db *sql.DB) CartRepo {
	return &cart{
		db: db,
	}
}

func NewOrder(db *sql.DB) OrderRepo {
	return &orders{
		db: db,
	}
}

func (p *pizzas) CreatePizzaType(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error) {

	query := `INSERT INTO types(name) 
	VALUES($1)`

	_, err := p.db.Exec(
		query,
		req.Name,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.CreatePizzaResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) CreatePizza(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error) {

	query := `INSERT INTO pizza(name, cost, type_id) 
	VALUES($1, $2, $3)`

	_, err := p.db.Exec(
		query,
		req.Name,
		req.Price,
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

	query := `SELECT name, cost FROM pizza WHERE id = $1 AND type_id= $2`

	err := p.db.QueryRow(query, req.Id, req.TypeId).Scan(&pizzas.Name, &pizzas.Price)
	if err != nil {
		return nil, err
	}

	return &pizza.GetPizzaByIdResponse{
		Name:  pizzas.Name,
		Price: pizzas.Price,
	}, nil
}

func (p *pizzas) GetPizzas(ctx context.Context, req *pizza.GetPizzasRequest) (*pizza.GetPizzasResponse, error) {

	query := `SELECT name, cost FROM pizza`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pizzas []*pizza.Pizzas

	for rows.Next() {
		var name string
		var price float32
		if err := rows.Scan(&name, &price); err != nil {
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

	query := `UPDATE pizza SET name = $3, cost = $4 WHERE id = $1 AND type_id = $2`

	_, err := p.db.Exec(
		query,
		req.Id,
		req.TypeId,
		req.Name,
		req.Price,
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

	query := `DELETE FROM pizza WHERE id = $1 AND type_id = $2`

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
		fmt.Println("error", err)
		return nil, sql.ErrNoRows
	} else if err != nil {
		fmt.Println("error", err)
		return nil, err
	}

	return &pizza.DeletePizzaResponse{
		Message: "success",
	}, nil
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

func (p *pizzas) GetPizzaCost(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	var items models.CartItems

	query := `
    	SELECT cost FROM pizza WHERE id = $1
	`

	err := p.db.QueryRow(query, req.PizzaId).Scan(&items.Cost)
	if err != nil {
		return nil, err
	}

	return &pizza.CartItemsResp{
		Cost: items.Cost,
	}, nil
}

func (c *cart) GetFromCart(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	query := `
    SELECT pizza_id, cart_id FROM cart_item WHERE id = $1
`
	// Tip: it's safer to use temporary variables with Scan instead of scanning directly into request struct (req).
	var pizzaId, cartId int32
	err := c.db.QueryRow(query, req.Id).Scan(&pizzaId, &cartId)
	if err != nil {
		return nil, err
	}

	return &pizza.CartItemsResp{
		CartId:  cartId,
		PizzaId: pizzaId,
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
    	SELECT COALESCE(SUM(ci.cost * ci.quantity), 0)
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
	UPDATE cart_item SET id = $1 WHERE id = $1
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

func (c *cart) DecreaseTotalCost(ctx context.Context, userId int32) (*pizza.CartItemsResp, error) {
	query := `
	UPDATE cart
	SET id = $1
	WHERE user_id = $1;
	`
	_, err := c.db.Exec(
		query,
		userId,
	)
	if err != nil {
		return nil, err
	}

	return &pizza.CartItemsResp{
		Message: "success",
	}, nil
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

// General tip: Always try to keep your code clean from debug prints in production, for local development you can use a logger.
// Keep variable and method names typo-free (typo-free means no spelling mistakes), and keep an eye on SQL query parameters and Scan argument counts for reliability.
