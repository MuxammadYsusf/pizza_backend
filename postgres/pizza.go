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
}

func NewPizza(db *sql.DB) PizzaRepo {
	return &pizzas{
		db: db,
	}
}

func (p *pizzas) CreatePizza(ctx context.Context, req *pizza.CreatePizzaRequest) (*pizza.CreatePizzaResponse, error) {

	query := `INSERT INTO pizza VALUES(name, price, id, typeId)`

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

	query := `SELECT name, price FROM pizza WHERE id = $1 AND typeId = $2`

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

	query := `UPDATE pizza SET name = $3, price = $4 WHERE id = $1 AND typeId = $2`

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
