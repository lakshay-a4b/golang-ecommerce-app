// models/order.go
package models

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Order struct {
	OrderID     int             `json:"id"`
	PaymentID   string          `json:"paymentId"`
	UserId      string          `json:"userId"`
	ProductInfo json.RawMessage `json:"productInfo"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"createdAt"`
}

type Cart struct {
	UserID      string          `json:"userId"`
	ProductInfo json.RawMessage `json:"productInfo"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type OrderRepository struct {
	pool *pgxpool.Pool
}

type CartRepository struct {
	pool *pgxpool.Pool
}

type OrderItem struct {
    ProductID int     `json:"productId"`
    Name      string  `json:"name"`
    Image     string  `json:"image"`
    Price     float64 `json:"price"`
    Quantity  int     `json:"quantity"`
}
