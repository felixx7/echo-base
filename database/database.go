package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"echo-base/config"
)

// Connect creates a new database connection
func Connect(dbConfig *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConfig.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// Close closes the database connection
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
