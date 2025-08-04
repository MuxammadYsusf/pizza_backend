package postgres

import "database/sql"

type NewPostgresI interface {
	Pizza() PizzaRepo
	Login() AuthRepo
}

type store struct {
	pizza PizzaRepo
	auth  AuthRepo
	db    *sql.DB
}

func (s *store) Pizza() PizzaRepo {
	return s.pizza
}

func (s *store) Login() AuthRepo {
	return s.auth
}

func NewPostgres(db *sql.DB) NewPostgresI {
	return &store{
		db:    db,
		pizza: NewPizza(db),
		auth:  NewAuth(db),
	}
}
