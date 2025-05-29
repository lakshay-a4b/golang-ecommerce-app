package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB *pgxpool.Pool
	// Add other configuration variables if needed
)

// InitDB initializes the database connection pool
func InitDB() error {
	// Validate required environment variables
	requiredVars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("missing required environment variable: %s", v)
		}
	}

	// Construct DSN with additional connection parameters
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable pool_max_conns=10 pool_min_conns=2 pool_max_conn_lifetime=1h",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Configure connection pool
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	// Set additional connection pool settings
	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	// Create connection pool with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	err = DB.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("PostgreSQL connection pool established successfully")
	return nil
}

// CloseDB closes the database connection pool
func CloseDB() error {
	if DB != nil {
		log.Println("Closing database connection pool...")
		DB.Close()
	}
	return nil
}

// GetDB returns the database connection pool
func GetDB() *pgxpool.Pool {
	return DB
}