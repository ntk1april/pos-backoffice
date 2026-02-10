package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/sijms/go-ora/v2"
	"pos-backoffice/internal/config"
)

var DB *sql.DB

// InitDB initializes Oracle database connection
func InitDB() error {
	dsn := config.AppConfig.GetDSN()
	
	var err error
	DB, err = sql.Open("oracle", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(10 * time.Minute)

	// Test connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✓ Database connection established")
	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		log.Println("Closing database connection...")
		return DB.Close()
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
}
