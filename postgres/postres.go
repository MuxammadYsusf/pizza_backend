package postgres

import "database/sql"

type NewPostgresI interface {
	Pizza() PizzaRepo
	Cart() CartRepo
	Order() OrderRepo
	Auth() AuthRepo
}

type store struct {
	pizza PizzaRepo
	order OrderRepo
	cart  CartRepo
	auth  AuthRepo
	db    *sql.DB
}

func (s *store) Pizza() PizzaRepo {
	return s.pizza
}

func (s *store) Cart() CartRepo {
	return s.cart
}

func (s *store) Order() OrderRepo {
	return s.order
}

func (s *store) Auth() AuthRepo {
	return s.auth
}

func NewPostgres(db *sql.DB) NewPostgresI {
	return &store{
		db:    db,
		pizza: NewPizza(db),
		order: NewOrder(db),
		cart:  NewCart(db),
		auth:  NewAuth(db),
	}
}
