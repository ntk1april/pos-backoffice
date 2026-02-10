package service

import (
	"database/sql"
	"fmt"

	"pos-backoffice/internal/models"
	"pos-backoffice/internal/repository"
)

type StockService struct {
	db          *sql.DB
	productRepo *repository.ProductRepository
	stockRepo   *repository.StockRepository
}

func NewStockService(db *sql.DB, productRepo *repository.ProductRepository, stockRepo *repository.StockRepository) *StockService {
	return &StockService{
		db:          db,
		productRepo: productRepo,
		stockRepo:   stockRepo,
	}
}

// IncreaseStock increases product stock with transaction
func (s *StockService) IncreaseStock(req *models.StockAdjustmentRequest, userID int64) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock product row for update
	product, err := s.productRepo.FindByIDForUpdate(tx, req.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.Status != "ACTIVE" {
		return fmt.Errorf("cannot adjust stock for inactive product")
	}

	// Calculate new stock
	stockBefore := product.Stock
	stockAfter := stockBefore + req.Quantity

	// Update stock
	err = s.productRepo.UpdateStock(tx, req.ProductID, stockAfter)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Create stock log
	stockLog := &models.StockLog{
		ProductID:       req.ProductID,
		TransactionType: "INCREASE",
		Quantity:        req.Quantity,
		StockBefore:     stockBefore,
		StockAfter:      stockAfter,
		Notes:           req.Notes,
		CreatedBy:       userID,
	}

	err = s.stockRepo.CreateLog(tx, stockLog)
	if err != nil {
		return fmt.Errorf("failed to create stock log: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DecreaseStock decreases product stock with transaction
func (s *StockService) DecreaseStock(req *models.StockAdjustmentRequest, userID int64) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock product row for update
	product, err := s.productRepo.FindByIDForUpdate(tx, req.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.Status != "ACTIVE" {
		return fmt.Errorf("cannot adjust stock for inactive product")
	}

	// Check if sufficient stock
	if product.Stock < req.Quantity {
		return fmt.Errorf("insufficient stock: available %d, requested %d", product.Stock, req.Quantity)
	}

	// Calculate new stock
	stockBefore := product.Stock
	stockAfter := stockBefore - req.Quantity

	// Update stock
	err = s.productRepo.UpdateStock(tx, req.ProductID, stockAfter)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Create stock log
	stockLog := &models.StockLog{
		ProductID:       req.ProductID,
		TransactionType: "DECREASE",
		Quantity:        req.Quantity,
		StockBefore:     stockBefore,
		StockAfter:      stockAfter,
		Notes:           req.Notes,
		CreatedBy:       userID,
	}

	err = s.stockRepo.CreateLog(tx, stockLog)
	if err != nil {
		return fmt.Errorf("failed to create stock log: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetStockLogs retrieves stock logs for a product
func (s *StockService) GetStockLogs(productID int64, page, pageSize int) (*models.StockLogResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := s.stockRepo.FindByProductID(productID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock logs: %w", err)
	}

	return &models.StockLogResponse{
		Logs:     logs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
