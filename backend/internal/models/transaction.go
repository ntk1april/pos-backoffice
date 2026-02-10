package models

import "time"

// Transaction represents a stock movement (INCREASE or DECREASE)
type Transaction struct {
	ID              int64     `json:"id"`
	TransactionType string    `json:"transaction_type"` // INCREASE or DECREASE
	ProductID       int64     `json:"product_id"`
	ProductName     string    `json:"product_name,omitempty"`
	StoreID         *int64    `json:"store_id"`         // NULL for INCREASE, NOT NULL for DECREASE
	StoreName       string    `json:"store_name,omitempty"`
	Quantity        int       `json:"quantity"`
	UnitPrice       float64   `json:"unit_price"`
	TotalAmount     float64   `json:"total_amount"`
	Notes           string    `json:"notes"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedBy       int64     `json:"created_by"`
	CreatedByName   string    `json:"created_by_name,omitempty"`
}

// TransactionRequest for creating new transactions
type TransactionRequest struct {
	TransactionType string  `json:"transaction_type" binding:"required,oneof=INCREASE DECREASE"`
	ProductID       int64   `json:"product_id" binding:"required"`
	StoreID         *int64  `json:"store_id"` // Required for DECREASE, NULL for INCREASE
	Quantity        int     `json:"quantity" binding:"required,min=1"`
	UnitPrice       float64 `json:"unit_price" binding:"required,min=0"`
	Notes           string  `json:"notes"`
}
