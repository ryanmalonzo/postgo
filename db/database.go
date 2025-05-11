package db

import (
	"fmt"

	_ "github.com/lib/pq"
)

// The CreateDatabase method creates a new PostgreSQL database using a Connection instance.
func (c *Connection) CreateDatabase(dbname string) error {
	// Check if the database already exists
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	err := c.db.QueryRow(query, dbname).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// If database doesn't exist, create it
	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE %s", dbname)
		_, err = c.db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database %s: %w", dbname, err)
		}
	}

	return nil
}
