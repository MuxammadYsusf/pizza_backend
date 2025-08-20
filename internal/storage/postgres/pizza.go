package postgres

import (
	"context"
	"database/sql"
	pb "github/http/copy/task4/genproto"
	"github/http/copy/task4/models"
)

type pizzas struct {
	db *sql.DB
}

type PizzaRepo interface {
	CreatePizzaType(ctx context.Context, req *pb.CreatePizzaRequest) (*pb.CreatePizzaResponse, error)
	CreatePizza(ctx context.Context, req *pb.CreatePizzaRequest) (*pb.CreatePizzaResponse, error)
	GetPizzaById(ctx context.Context, req *pb.GetPizzaByIdRequest) (*pb.GetPizzaByIdResponse, error)
	GetPizzas(ctx context.Context, req *pb.GetPizzasRequest) (*pb.GetPizzasResponse, error)
	UpdatePizza(ctx context.Context, req *pb.UpdatePizzaRequest) (*pb.UpdatePizzaResponse, error)
	DeletePizza(ctx context.Context, req *pb.DeletePizzaRequest) (*pb.DeletePizzaResponse, error)
	GetPizzaData(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error)
	GetPizzaCost(ctx context.Context, pizzaId int32) (*pb.CartItemsResp, error)
	GetAllPizzaCost(ctx context.Context, orderId int32) (*pb.OrderPizzaResponse, error)
}

func NewPizza(db *sql.DB) PizzaRepo {
	return &pizzas{
		db: db,
	}
}

func (p *pizzas) CreatePizzaType(ctx context.Context, req *pb.CreatePizzaRequest) (*pb.CreatePizzaResponse, error) {

	query := `INSERT INTO types(name) 
	VALUES($1)`

	_, err := p.db.Exec(
		query,
		req.Name,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePizzaResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) CreatePizza(ctx context.Context, req *pb.CreatePizzaRequest) (*pb.CreatePizzaResponse, error) {

	query := `INSERT INTO pizza(name, cost, type_id, photo) 
	VALUES($1, $2, $3, $4)`

	_, err := p.db.Exec(
		query,
		req.Name,
		req.Price,
		req.TypeId,
		req.Photo,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePizzaResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) GetPizzaById(ctx context.Context, req *pb.GetPizzaByIdRequest) (*pb.GetPizzaByIdResponse, error) {
	var pizzas models.Pizza

	query := `SELECT name, cost, photo FROM pizza WHERE id = $1 AND type_id= $2`

	err := p.db.QueryRow(query, req.Id, req.TypeId).Scan(&pizzas.Name, &pizzas.Price, &pizzas.Photo)
	if err != nil {
		return nil, err
	}

	return &pb.GetPizzaByIdResponse{
		Name:  pizzas.Name,
		Price: pizzas.Price,
		Photo: pizzas.Photo,
	}, nil
}

func (p *pizzas) GetPizzas(ctx context.Context, req *pb.GetPizzasRequest) (*pb.GetPizzasResponse, error) {

	query := `SELECT id, name, cost, photo FROM pizza ORDER BY id ASC`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pizzas []*pb.Pizzas

	for rows.Next() {
		var p models.Pizza
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Photo); err != nil {
			return nil, err
		}

		pizzas = append(pizzas, &pb.Pizzas{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
			Photo: p.Photo,
		})
	}

	return &pb.GetPizzasResponse{
		Pizzas: pizzas,
	}, nil
}

func (p *pizzas) UpdatePizza(ctx context.Context, req *pb.UpdatePizzaRequest) (*pb.UpdatePizzaResponse, error) {

	query := `UPDATE pizza SET name = $3, cost = $4, photo = $5 WHERE id = $1 AND type_id = $2`

	_, err := p.db.Exec(
		query,
		req.Id,
		req.TypeId,
		req.Name,
		req.Price,
		req.Photo,
	)

	if err != nil {
		return nil, err
	}

	return &pb.UpdatePizzaResponse{
		Message: "success",
		Name:    req.Name,
		Price:   req.Price,
	}, nil
}

func (p *pizzas) DeletePizza(ctx context.Context, req *pb.DeletePizzaRequest) (*pb.DeletePizzaResponse, error) {

	query := `DELETE FROM pizza WHERE id = $1`

	result, err := p.db.Exec(
		query,
		req.Id,
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

	return &pb.DeletePizzaResponse{
		Message: "success",
	}, nil
}

func (p *pizzas) GetPizzaData(ctx context.Context, req *pb.CartItems) (*pb.CartItemsResp, error) {

	var items models.CartItems

	query := `
    	SELECT cost, id FROM pizza WHERE id = $1
	`

	err := p.db.QueryRow(query, req.PizzaId).Scan(&items.Cost, &items.ID)
	if err != nil {
		return nil, err
	}

	return &pb.CartItemsResp{
		Cost:    items.Cost,
		PizzaId: items.ID,
	}, nil
}

func (p *pizzas) GetPizzaCost(ctx context.Context, pizzaId int32) (*pb.CartItemsResp, error) {
	var cost float32

	query := `SELECT COALESCE(SUM(p.cost * ci.quantity), 0) AS cost
	FROM pizza p
	JOIN cart_item ci ON p.id = ci.pizza_id
	WHERE p.id = $1`

	err := p.db.QueryRow(query, pizzaId).Scan(&cost)
	if err != nil {
		return nil, err
	}

	return &pb.CartItemsResp{
		Cost: cost,
	}, nil
}

func (p *pizzas) GetAllPizzaCost(ctx context.Context, orderId int32) (*pb.OrderPizzaResponse, error) {
	var costs []float32

	query := `SELECT 
	    ci.pizza_id,
	    ci.quantity,
	    (p.cost * ci.quantity) AS cost
	FROM cart_item ci
	JOIN pizza p ON p.id = ci.pizza_id
	WHERE ci.cart_id = $1;`

	rows, err := p.db.Query(query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pizzaId, quantity int32
		var cost float64
		if err := rows.Scan(&pizzaId, &quantity, &cost); err != nil {
			return nil, err
		}
		costs = append(costs, float32(cost))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.OrderPizzaResponse{
		Cost: costs,
	}, nil
}
