package models

import "time"

type Payment struct {
	PaymentID   string  `json:"paymentId"`
	UserId      string  `json:"userId"`
	TotalAmount float64 `json:"totalAmount"`
	Status      string  `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type PaymentResponse struct {
    TransactionID string
    Status        string
}


type PaymentRequest struct {
    UserID string
    Amount float64
}