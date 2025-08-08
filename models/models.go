package models

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Pizza struct {
	ID       int32   `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	TypeId   int32   `json:"typeId"`
	Quantity int32   `json:"quantity"`
}

type Cart struct {
	ID        int     `json:"id"`
	UserId    int     `json:"userId"`
	IsActive  bool    `json:"isActive"`
	TotalCost float32 `json:"totalCost"`
}

type CartItems struct {
	ID          int32   `json:"id"`
	PizzaId     int32   `json:"pizzaId"`
	PizzaTypeId int32   `json:"pizzaTypeId"`
	Cost        float32 `json:"cost"`
	CartId      int32   `json:"cartId"`
	Quantity    int32   `json:"quantity"`
	TotalCost   float32 `json:"totalCost"`
}

type CartIeamHistory struct {
	ID          int32   `json:"id"`
	PizzaId     int32   `json:"pizzaId"`
	PizzaTypeId int32   `json:"pizzaTypeId"`
	Cost        float32 `json:"cost"`
	Quantity    int32   `json:"quantity"`
	TotalCost   float32 `json:"totalCost"`
}

type Order struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	IsOrdered bool      `json:"isOrdered"`
	UserId    int       `json:"userId"`
	Status    string    `json:"status"`
}
