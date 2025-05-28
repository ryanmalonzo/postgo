package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// CreateDatabase crée une nouvelle base de données PostgreSQL en utilisant une instance Connection.
// Cette fonction vérifie d'abord si la base de données existe déjà avant de la créer
// pour éviter les erreurs de duplication.
func (c *Connection) CreateDatabase(dbname string) error {
	// Vérification si la base de données existe déjà
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	err := c.db.QueryRow(query, dbname).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Si la base de données n'existe pas, la créer
	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE %s", dbname)
		_, err = c.db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database %s: %w", dbname, err)
		}
	}

	return nil
}
func (c *Connection) GetDatabase() *sql.DB {
	return c.db
}
