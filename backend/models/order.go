package models

import "time"

// order represents a transaction in the Ledger(postgress)

type Order struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	TotalAmountCents int64     `json:"total_amount_cents"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at`
}

type CreateOrderRequest struct {
	ProductID int64 `json:"product_id"`
	// in real app userid would come from JWT.
	UserID string `json:"user_id"`
}
