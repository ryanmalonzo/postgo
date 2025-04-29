package functions

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(host string, port int, user, password, dbname string) (*sql.DB, error) {
	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
