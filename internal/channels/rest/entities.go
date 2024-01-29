package rest

import (
	"time"
)

type Response struct {
	Message string `json:"message"`
}

type PaymentRequest struct {
	PaymentType int       `json:"payment_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      int       `json:"status"`
	OrderID     string    `json:"order_id"`
}

type PaymentCallback struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}
