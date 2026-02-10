package repository

import (
	"database/sql"
	"fmt"

	"pos-backoffice/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByUsername finds a user by username using raw SQL
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, full_name, role, status, created_at, updated_at
		FROM users
		WHERE username = :1 AND status = 'ACTIVE'
	`

	var user models.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, full_name, role, status, created_at, updated_at
		FROM users
		WHERE id = :1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, password_hash, full_name, role, status)
		VALUES (:1, :2, :3, :4, :5)
		RETURNING id INTO :6
	`

	_, err := r.db.Exec(query,
		user.Username,
		user.PasswordHash,
		user.FullName,
		user.Role,
		user.Status,
		sql.Out{Dest: &user.ID},
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
