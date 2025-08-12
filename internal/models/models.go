package models

import "time"

// User представляет пользователя в системе.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Product представляет товар (пряжа или готовое изделие).
type Product struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Image           string  `json:"image"`
	CategoryID      int     `json:"category_id"`
	Type            string  `json:"type"`
	Composition     string  `json:"composition,omitempty"`
	CountryOfOrigin string  `json:"country_of_origin,omitempty"`
	LengthIn100g    int     `json:"length_in_100g,omitempty"`
	Size            string  `json:"size,omitempty"`
	GarmentLength   string  `json:"garment_length,omitempty"`
	Color           string  `json:"color,omitempty"`
}

// Category представляет категорию товаров.
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Order представляет заказ, сделанный пользователем.
type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderItem представляет отдельный товар в заказе.
type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Comment представляет комментарий к товару.
type Comment struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	UserID    int       `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
