package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Connection struct {
	db *sql.DB
}

// NewConnection creates a new database connection using individual connection parameters
func NewConnection(host string, port int, user, password, dbname string) (*Connection, error) {
	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Connection{db: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.db.Close()
}
