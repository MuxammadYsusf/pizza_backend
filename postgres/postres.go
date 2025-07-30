package postgres

import "database/sql"

type NewPostgresI interface {
	Pizza() PizzaRepo
	Login() LoginRepo
}

type store struct {
	pizza PizzaRepo
	login LoginRepo
	db    *sql.DB
}

func (s *store) Pizza() PizzaRepo {
	return s.pizza
}

func (s *store) Login() LoginRepo {
	return s.login
}

func NewPostgres(db *sql.DB) NewPostgresI {
	return &store{
		db:    db,
		pizza: NewPizza(db),
		login: NewLogin(db),
	}
}
