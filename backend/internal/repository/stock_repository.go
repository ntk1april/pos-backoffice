package repository

import (
	"database/sql"
	"fmt"

	"pos-backoffice/internal/models"
)

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

// CreateLog creates a stock log entry within a transaction
func (r *StockRepository) CreateLog(tx *sql.Tx, log *models.StockLog) error {
	query := `
		INSERT INTO stock_logs (product_id, transaction_type, quantity, stock_before, stock_after, notes, created_by)
		VALUES (:1, :2, :3, :4, :5, :6, :7)
		RETURNING id INTO :8
	`

	_, err := tx.Exec(query,
		log.ProductID,
		log.TransactionType,
		log.Quantity,
		log.StockBefore,
		log.StockAfter,
		log.Notes,
		log.CreatedBy,
		sql.Out{Dest: &log.ID},
	)

	if err != nil {
		return fmt.Errorf("failed to create stock log: %w", err)
	}

	return nil
}

// FindByProductID retrieves stock logs for a product with pagination
func (r *StockRepository) FindByProductID(productID int64, page, pageSize int) ([]models.StockLog, int64, error) {
	offset := (page - 1) * pageSize

	// Count total records
	countQuery := "SELECT COUNT(*) FROM stock_logs WHERE product_id = :1"
	var total int64
	err := r.db.QueryRow(countQuery, productID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stock logs: %w", err)
	}

	// Query with pagination
	query := fmt.Sprintf(`
		SELECT id, product_id, transaction_type, quantity, stock_before, stock_after, notes, created_by, created_at
		FROM stock_logs
		WHERE product_id = :1
		ORDER BY created_at DESC
		OFFSET %d ROWS FETCH NEXT %d ROWS ONLY
	`, offset, pageSize)

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query stock logs: %w", err)
	}
	defer rows.Close()

	var logs []models.StockLog
	for rows.Next() {
		var log models.StockLog
		err := rows.Scan(
			&log.ID,
			&log.ProductID,
			&log.TransactionType,
			&log.Quantity,
			&log.StockBefore,
			&log.StockAfter,
			&log.Notes,
			&log.CreatedBy,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan stock log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

// GetRecentLogs retrieves recent stock logs across all products
func (r *StockRepository) GetRecentLogs(limit int) ([]models.StockLog, error) {
	query := fmt.Sprintf(`
		SELECT id, product_id, transaction_type, quantity, stock_before, stock_after, notes, created_by, created_at
		FROM stock_logs
		ORDER BY created_at DESC
		FETCH FIRST %d ROWS ONLY
	`, limit)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent stock logs: %w", err)
	}
	defer rows.Close()

	var logs []models.StockLog
	for rows.Next() {
		var log models.StockLog
		err := rows.Scan(
			&log.ID,
			&log.ProductID,
			&log.TransactionType,
			&log.Quantity,
			&log.StockBefore,
			&log.StockAfter,
			&log.Notes,
			&log.CreatedBy,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
