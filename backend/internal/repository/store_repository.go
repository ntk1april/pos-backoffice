package repository

import (
	"database/sql"
	"fmt"
	"pos-backoffice/internal/models"
)

type StoreRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// GetAll returns all stores
func (r *StoreRepository) GetAll() ([]models.Store, error) {
	query := `
		SELECT id, code, name, address, phone, status, 
		       created_at, updated_at, created_by, updated_by
		FROM stores
		WHERE status = 'ACTIVE'
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []models.Store
	for rows.Next() {
		var store models.Store
		err := rows.Scan(
			&store.ID, &store.Code, &store.Name, &store.Address, &store.Phone,
			&store.Status, &store.CreatedAt, &store.UpdatedAt,
			&store.CreatedBy, &store.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	return stores, nil
}

// GetByID returns a store by ID
func (r *StoreRepository) GetByID(id int64) (*models.Store, error) {
	query := `
		SELECT id, code, name, address, phone, status, 
		       created_at, updated_at, created_by, updated_by
		FROM stores
		WHERE id = :1
	`

	var store models.Store
	err := r.db.QueryRow(query, id).Scan(
		&store.ID, &store.Code, &store.Name, &store.Address, &store.Phone,
		&store.Status, &store.CreatedAt, &store.UpdatedAt,
		&store.CreatedBy, &store.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &store, nil
}

// Create creates a new store
func (r *StoreRepository) Create(store *models.Store) error {
	query := `
		INSERT INTO stores (code, name, address, phone, status, created_by, updated_by)
		VALUES (:1, :2, :3, :4, :5, :6, :7)
		RETURNING id INTO :8
	`

	_, err := r.db.Exec(query,
		store.Code, store.Name, store.Address, store.Phone,
		store.Status, store.CreatedBy, store.UpdatedBy,
		sql.Out{Dest: &store.ID},
	)

	return err
}

// Update updates an existing store
func (r *StoreRepository) Update(store *models.Store) error {
	query := `
		UPDATE stores
		SET code = :1, name = :2, address = :3, phone = :4, 
		    status = :5, updated_by = :6, updated_at = CURRENT_TIMESTAMP
		WHERE id = :7
	`

	result, err := r.db.Exec(query,
		store.Code, store.Name, store.Address, store.Phone,
		store.Status, store.UpdatedBy, store.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("store not found")
	}

	return nil
}

// Delete soft deletes a store
func (r *StoreRepository) Delete(id int64) error {
	query := `UPDATE stores SET status = 'INACTIVE' WHERE id = :1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("store not found")
	}

	return nil
}
