package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github/http/copy/task4/generated/pizza"
	"github/http/copy/task4/models"
)

type pizzas struct {
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
