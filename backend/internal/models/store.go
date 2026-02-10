package models

import "time"

type Store struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
}

type StoreRequest struct {
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Status  string `json:"status"`
}
