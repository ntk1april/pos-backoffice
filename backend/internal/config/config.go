package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBService   string
	DBUser      string
	DBPassword  string
	JWTSecret   string
	ServerPort  string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() error {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "1521"),
		DBService:  getEnv("DB_SERVICE", "XEPDB1"),
		DBUser:     getEnv("DB_USER", "pos_user"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		ServerPort: getEnv("PORT", "8080"),
	}

	// Validate required fields
	if AppConfig.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	if AppConfig.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	return nil
}

// GetDSN returns Oracle connection string for go-ora
func (c *Config) GetDSN() string {
	// go-ora format: oracle://user:password@host:port/service_name
	return fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBService,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
