package repository

import (
	"database/sql"
	"pos-backoffice/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction and updates product stock
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	// Start database transaction
	dbTx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer dbTx.Rollback()

	// Insert transaction
	query := `
		INSERT INTO transactions (
			transaction_type, product_id, store_id, quantity, 
			unit_price, total_amount, notes, created_by
		)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8)
		RETURNING id, transaction_date INTO :9, :10
	`

	_, err = dbTx.Exec(query,
		tx.TransactionType, tx.ProductID, tx.StoreID, tx.Quantity,
		tx.UnitPrice, tx.TotalAmount, tx.Notes, tx.CreatedBy,
		sql.Out{Dest: &tx.ID},
		sql.Out{Dest: &tx.TransactionDate},
	)

	if err != nil {
		return err
	}

	// Update product stock
	var stockUpdate string
	if tx.TransactionType == "INCREASE" {
		stockUpdate = `UPDATE products SET stock = stock + :1, updated_at = CURRENT_TIMESTAMP WHERE id = :2`
	} else {
		stockUpdate = `UPDATE products SET stock = stock - :1, updated_at = CURRENT_TIMESTAMP WHERE id = :2`
	}

	_, err = dbTx.Exec(stockUpdate, tx.Quantity, tx.ProductID)
	if err != nil {
		return err
	}

	return dbTx.Commit()
}

// GetByProductID returns all transactions for a product
func (r *TransactionRepository) GetByProductID(productID int64, limit int) ([]models.Transaction, error) {
	query := `
		SELECT 
			t.id, t.transaction_type, t.product_id, p.name as product_name,
			t.store_id, s.name as store_name,
			t.quantity, t.unit_price, t.total_amount, t.notes,
			t.transaction_date, t.created_by, u.full_name as created_by_name
		FROM transactions t
		JOIN products p ON t.product_id = p.id
		LEFT JOIN stores s ON t.store_id = s.id
		JOIN users u ON t.created_by = u.id
		WHERE t.product_id = :1
		ORDER BY t.transaction_date DESC
		FETCH FIRST :2 ROWS ONLY
	`

	rows, err := r.db.Query(query, productID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		var storeID sql.NullInt64
		var storeName sql.NullString
		var notes sql.NullString

		err := rows.Scan(
			&tx.ID, &tx.TransactionType, &tx.ProductID, &tx.ProductName,
			&storeID, &storeName,
			&tx.Quantity, &tx.UnitPrice, &tx.TotalAmount, &notes,
			&tx.TransactionDate, &tx.CreatedBy, &tx.CreatedByName,
		)

		if err != nil {
			return nil, err
		}

		if storeID.Valid {
			tx.StoreID = &storeID.Int64
		}

		if storeName.Valid {
			tx.StoreName = storeName.String
		}

		if notes.Valid {
			tx.Notes = notes.String
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// GetAll returns all transactions with pagination
func (r *TransactionRepository) GetAll(page, limit int) ([]models.Transaction, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM transactions`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get transactions
	query := `
		SELECT 
			t.id, t.transaction_type, t.product_id, p.name as product_name,
			t.store_id, s.name as store_name,
			t.quantity, t.unit_price, t.total_amount, t.notes,
			t.transaction_date, t.created_by, u.full_name as created_by_name
		FROM transactions t
		JOIN products p ON t.product_id = p.id
		LEFT JOIN stores s ON t.store_id = s.id
		JOIN users u ON t.created_by = u.id
		ORDER BY t.transaction_date DESC
		OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY
	`

	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		var storeID sql.NullInt64
		var storeName sql.NullString
		var notes sql.NullString

		err := rows.Scan(
			&tx.ID, &tx.TransactionType, &tx.ProductID, &tx.ProductName,
			&storeID, &storeName,
			&tx.Quantity, &tx.UnitPrice, &tx.TotalAmount, &notes,
			&tx.TransactionDate, &tx.CreatedBy, &tx.CreatedByName,
		)

		if err != nil {
			return nil, 0, err
		}

		if storeID.Valid {
			tx.StoreID = &storeID.Int64
		}

		if storeName.Valid {
			tx.StoreName = storeName.String
		}

		if notes.Valid {
			tx.Notes = notes.String
		}

		transactions = append(transactions, tx)
	}

	return transactions, total, nil
}

// GetByStoreID returns all transactions for a store
func (r *TransactionRepository) GetByStoreID(storeID int64, limit int) ([]models.Transaction, error) {
	query := `
		SELECT 
			t.id, t.transaction_type, t.product_id, p.name as product_name,
			t.store_id, s.name as store_name,
			t.quantity, t.unit_price, t.total_amount, t.notes,
			t.transaction_date, t.created_by, u.full_name as created_by_name
		FROM transactions t
		JOIN products p ON t.product_id = p.id
		JOIN stores s ON t.store_id = s.id
		JOIN users u ON t.created_by = u.id
		WHERE t.store_id = :1
		ORDER BY t.transaction_date DESC
		FETCH FIRST :2 ROWS ONLY
	`

	rows, err := r.db.Query(query, storeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		var notes sql.NullString

		err := rows.Scan(
			&tx.ID, &tx.TransactionType, &tx.ProductID, &tx.ProductName,
			&tx.StoreID, &tx.StoreName,
			&tx.Quantity, &tx.UnitPrice, &tx.TotalAmount, &notes,
			&tx.TransactionDate, &tx.CreatedBy, &tx.CreatedByName,
		)

		if err != nil {
			return nil, err
		}

		if notes.Valid {
			tx.Notes = notes.String
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
