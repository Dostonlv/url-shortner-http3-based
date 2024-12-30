// pkg/database/postgres.go
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates new PostgreSQL connection
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	// PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// Create tables if they don't exist
func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(10) NOT NULL UNIQUE,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		clicks INTEGER NOT NULL DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	return nil
}

// CloseDB safely closes database connection
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
