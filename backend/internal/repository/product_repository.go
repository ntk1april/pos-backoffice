package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"pos-backoffice/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// FindAll retrieves products with pagination and search
func (r *ProductRepository) FindAll(page, pageSize int, search string, status string) ([]models.Product, int64, error) {
	offset := (page - 1) * pageSize

	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		whereClause += fmt.Sprintf(" AND status = :%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if search != "" {
		whereClause += fmt.Sprintf(" AND UPPER(name) LIKE :%d", argIndex)
		args = append(args, "%"+strings.ToUpper(search)+"%")
		argIndex++
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Query with pagination using OFFSET/FETCH
	query := fmt.Sprintf(`
		SELECT id, sku, name, description, price, cost, stock, status, 
		       created_at, updated_at, created_by, updated_by
		FROM products
		%s
		ORDER BY created_at DESC
		OFFSET %d ROWS FETCH NEXT %d ROWS ONLY
	`, whereClause, offset, pageSize)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.SKU,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Cost,
			&p.Stock,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CreatedBy,
			&p.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, total, nil
}

// FindByID retrieves a product by ID
func (r *ProductRepository) FindByID(id int64) (*models.Product, error) {
	query := `
		SELECT id, sku, name, description, price, cost, stock, status,
		       created_at, updated_at, created_by, updated_by
		FROM products
		WHERE id = :1
	`

	var p models.Product
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.SKU,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Cost,
		&p.Stock,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.CreatedBy,
		&p.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query product: %w", err)
	}

	return &p, nil
}

// FindBySKU retrieves a product by SKU
func (r *ProductRepository) FindBySKU(sku string) (*models.Product, error) {
	query := `
		SELECT id, sku, name, description, price, cost, stock, status,
		       created_at, updated_at, created_by, updated_by
		FROM products
		WHERE sku = :1
	`

	var p models.Product
	err := r.db.QueryRow(query, sku).Scan(
		&p.ID,
		&p.SKU,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Cost,
		&p.Stock,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.CreatedBy,
		&p.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, nil // SKU not found is not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query product by SKU: %w", err)
	}

	return &p, nil
}

// Create creates a new product
func (r *ProductRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)
		RETURNING id INTO :10
	`

	_, err := r.db.Exec(query,
		product.SKU,
		product.Name,
		product.Description,
		product.Price,
		product.Cost,
		product.Stock,
		product.Status,
		product.CreatedBy,
		product.UpdatedBy,
		sql.Out{Dest: &product.ID},
	)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// Update updates an existing product
func (r *ProductRepository) Update(product *models.Product) error {
	query := `
		UPDATE products
		SET name = :1, description = :2, price = :3, cost = :4, updated_by = :5
		WHERE id = :6
	`

	result, err := r.db.Exec(query,
		product.Name,
		product.Description,
		product.Price,
		product.Cost,
		product.UpdatedBy,
		product.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// Delete soft deletes a product (sets status to INACTIVE)
func (r *ProductRepository) Delete(id int64, userID int64) error {
	query := `
		UPDATE products
		SET status = 'INACTIVE', updated_by = :1
		WHERE id = :2
	`

	result, err := r.db.Exec(query, userID, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// UpdateStock updates product stock (used within transactions)
func (r *ProductRepository) UpdateStock(tx *sql.Tx, productID int64, newStock int) error {
	query := `
		UPDATE products
		SET stock = :1
		WHERE id = :2
	`

	result, err := tx.Exec(query, newStock, productID)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// FindByIDForUpdate retrieves a product with row lock (FOR UPDATE)
func (r *ProductRepository) FindByIDForUpdate(tx *sql.Tx, id int64) (*models.Product, error) {
	query := `
		SELECT id, sku, name, description, price, cost, stock, status,
		       created_at, updated_at, created_by, updated_by
		FROM products
		WHERE id = :1
		FOR UPDATE
	`

	var p models.Product
	err := tx.QueryRow(query, id).Scan(
		&p.ID,
		&p.SKU,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Cost,
		&p.Stock,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.CreatedBy,
		&p.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query product: %w", err)
	}

	return &p, nil
}
