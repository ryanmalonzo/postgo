package functions

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)

// CreateDatabase creates a new PostgreSQL database if it doesn't already exist.
// It requires connection to the 'postgres' system database to perform the operation.
func CreateDatabase(host string, port int, user, password, dbname string) error {
    // Connect to the postgres database first to be able to create a new database
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", 
        host, port, user, password)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to connect to postgres database: %w", err)
    }
    defer db.Close()

    // Check if the database already exists
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
    err = db.QueryRow(query, dbname).Scan(&exists)
    if err != nil {
        return fmt.Errorf("failed to check if database exists: %w", err)
    }

    // If database doesn't exist, create it
    if !exists {
        // PostgreSQL doesn't support parameterized CREATE DATABASE statements
        // So we need to safely format the query
        createQuery := fmt.Sprintf("CREATE DATABASE %s", dbname)
        _, err = db.Exec(createQuery)
        if err != nil {
            return fmt.Errorf("failed to create database %s: %w", dbname, err)
        }
        return nil
    }

    // Database already exists
    return nil
}