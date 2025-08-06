package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/models"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type pizzas struct {
	db *sql.DB
}

type PizzaRepo interface {
	CreatePizza(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error)
	CreatePizzaType(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error)
	GetPizzaById(ctx context.Context, req *pizza.GetPizzaByIdRequest) (*pizza.GetPizzaByIdResponse, error)
	GetPizzas(ctx context.Context, req *pizza.GetPizzasRequest) (*pizza.GetPizzasResponse, error)
	UpdatePizza(ctx context.Context, req *pizza.UpdatePizzaRequest) (*pizza.UpdatePizzaResponse, error)
	DeletePizza(ctx context.Context, req *pizza.DeletePizzaRequest) (*pizza.DeletePizzaResponse, error)
	CheckIsCartExist(ctx context.Context, req *pizza.CheckIsCartExistRequest) (*pizza.CheckIsCartExistResponse, error)
	Cart(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error)
	CartItems(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error)
	GetPizzaCost(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
	GetFromCart(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
	GetFromPizza(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
	UpdatePizzaInCart(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error)
	UpdateTotalCost(ctx context.Context, id int32) (*pizza.CartItemsResp, error)
	CheckIsOrdered(ctx context.Context, req *pizza.CheckIsOrderedRequest) (*pizza.CheckIsOrderedResponse, error)
	Order(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	OrderItem(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error)
	GetCartrHistory(ctx context.Context, req *pizza.GetCartHistoryRequest) (*pizza.GetCartHistoryResponse, error)
	GetCartItemHistory(ctx context.Context, req *pizza.GetCarItemtHistoryRequest) (*pizza.GetCarItemtHistoryResponse, error)
}

func NewPizza(db *sql.DB) PizzaRepo {
	return &pizzas{
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
	fmt.Printf("req : %+v\n\n", req)
	if err != nil {
		fmt.Println("error", err)
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

func (p *pizzas) CheckIsCartExist(ctx context.Context, req *pizza.CheckIsCartExistRequest) (*pizza.CheckIsCartExistResponse, error) {
	var cart models.Cart

	query := `SELECT is_active FROM cart WHERE id = $1 AND user_id = $2`

	err := p.db.QueryRow(query, req.Id, req.UserId).Scan(&cart.IsActive)
	if err != nil {
		fmt.Println("cart.IsActive -->", cart.IsActive)
		fmt.Println("id -->", req.Id)
		fmt.Println("userId -->", req.UserId)
		return nil, err
	}

	return &pizza.CheckIsCartExistResponse{
		Message:  "success",
		IsActive: cart.IsActive,
	}, nil
}

func (p *pizzas) Cart(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error) {

	query := `INSERT INTO cart(user_id, is_active, total_cost) 
	VALUES($1, $2, $3)`

	_, err := p.db.Exec(
		query,
		req.UserId,
		req.IsActive,
		req.TotalCost,
	)
	fmt.Printf("[UserId: %d] [IsActive: %t] [TotalCost: %f]\n", req.UserId, req.IsActive, req.TotalCost)
	if err != nil {
		return nil, err
	}

	return &pizza.CartResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) CartItems(ctx context.Context, req *pizza.CartRequest) (*pizza.CartResponse, error) {

	query := `INSERT INTO cart_item(pizza_id, pizza_type_id, cost, cart_id, quantity)
	VALUES($1, $2, $3, $4, $5)`

	_, err := p.db.Exec(
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
		fmt.Println("HIHIHIHA:", err)
		return nil, err
	}

	return &pizza.CartItemsResp{
		Cost: items.Cost,
	}, nil
}

func (p *pizzas) GetFromCart(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	query := `
    SELECT pizza_id, cart_id FROM cart_item WHERE id = $1
`
	err := p.db.QueryRow(query, req.Id).Scan(&req.PizzaId, &req.CartId)
	if err != nil {
		return nil, err
	}

	return &pizza.CartItemsResp{
		CartId:  req.CartId,
		PizzaId: req.PizzaId,
	}, nil
}

func (p *pizzas) GetFromPizza(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	query := `
    SELECT cost FROM pizza WHERE id = $1
`
	err := p.db.QueryRow(query, req.PizzaId).Scan(&req.Cost)
	if err != nil {
		return nil, err
	}

	return &pizza.CartItemsResp{
		Cost: req.Cost,
	}, nil
}

func (p *pizzas) UpdatePizzaInCart(ctx context.Context, req *pizza.CartItems) (*pizza.CartItemsResp, error) {

	query := `UPDATE cart_item SET quantity = $3, cost = $4 WHERE id = $1 AND piazza_id = $2`

	tx, err := p.db.Begin()
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

	return &pizza.CartItemsResp{
		Message: "success",
	}, nil
}

func (p *pizzas) UpdateTotalCost(ctx context.Context, id int32) (*pizza.CartItemsResp, error) {
	query := `
	UPDATE cart
	SET total_cost = (
    	SELECT COALESCE(SUM(ci.cost * ci.quantity), 0)
    	FROM cart_item ci
 		WHERE ci.cart_id = $1
	)
	WHERE id = $1;
	`
	_, err := p.db.Exec(
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

func (p *pizzas) CheckIsOrdered(ctx context.Context, req *pizza.CheckIsOrderedRequest) (*pizza.CheckIsOrderedResponse, error) {
	var cart models.Order

	query := `SELECT is_ordered, status FROM orders WHERE id = $1 AND user_id = $2`

	err := p.db.QueryRow(query, req.Id, req.UserId).Scan(&cart.IsOrdered, &cart.Status)
	if err != nil {
		return nil, err
	}

	return &pizza.CheckIsOrderedResponse{
		Message:   "success",
		IsOrdered: cart.IsOrdered,
	}, nil
}

func (p *pizzas) Order(ctx context.Context, req *pizza.OrderPizzaRequest) (*pizza.OrderPizzaResponse, error) {

	query := `INSERT INTO orders(date, is_ordered, user_id, status) 
	VALUES($1, $2, $3, $4)`

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

	query := `INSERT INTO order_item(pizza_id, pizza_type_id, cart_id, quantity) 
	VALUES($1, $2, $3, $4)`

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

func (p *pizzas) GetCartrHistory(ctx context.Context, req *pizza.GetCartHistoryRequest) (*pizza.GetCartHistoryResponse, error) {

	query := `SELECT c.id, c.is_active, o.date
          FROM cart AS c
          INNER JOIN orders AS o ON o.cart_id = c.id
          WHERE c.user_id = $1`

	rows, err := p.db.Query(query, req.CartId, req.UserId)
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

func (p *pizzas) GetCartItemHistory(ctx context.Context, req *pizza.GetCarItemtHistoryRequest) (*pizza.GetCarItemtHistoryResponse, error) {

	var cart models.CartIeamHistory

	query := `SELECT ci.pizza_id, ci.pizza_type_id, ci.cost, ci.quantity, c.total_cost
	FROM cart_item AS ci 
	JOIN cart AS c ON c.id = ci.cart_id
	WHERE c.id = $1`

	err := p.db.QueryRow(query, req.CartHistoryId).Scan(&cart.PizzaId, &cart.PizzaTypeId, &cart.Cost, &cart.Quantity)
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
