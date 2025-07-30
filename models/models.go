package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Pizza struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	TypeId   int     `json:"typeId"`
	Quantity int     `json:"quantity"`
}
