package models

import "time"

type User struct {
	ID       int32  `json:"id"`
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
	Photo    string  `json:"photo"`
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
	ID        int32     `json:"id"`
	Date      time.Time `json:"date"`
	IsOrdered bool      `json:"isOrdered"`
	UserId    int32     `json:"userId"`
	Status    string    `json:"status"`
}

type OrderItem struct {
	Id       int32   `json:"id"`
	PizzaId  int32   `json:"pizzaId"`
	Quantity int32   `json:"quantity"`
	OrderId  int32   `json:"orderId"`
	Caot     float32 `json:"cost"`
}

type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiredAt time.Time `json:"expiresAt"`
}
