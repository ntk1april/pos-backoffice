package models

import "time"

type StockLog struct {
	ID              int64     `json:"id"`
	ProductID       int64     `json:"product_id"`
	TransactionType string    `json:"transaction_type"` // INCREASE, DECREASE, ADJUSTMENT
	Quantity        int       `json:"quantity"`
	StockBefore     int       `json:"stock_before"`
	StockAfter      int       `json:"stock_after"`
	Notes           string    `json:"notes"`
	CreatedBy       int64     `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}

type StockAdjustmentRequest struct {
	ProductID int64  `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
	Notes     string `json:"notes"`
}

type StockLogResponse struct {
	Logs      []StockLog `json:"logs"`
	Total     int64      `json:"total"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
}
