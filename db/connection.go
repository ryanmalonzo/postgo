package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Connection struct {
	db *sql.DB
}

// NewConnection initialise une nouvelle connexion à la base de données PostgreSQL.
// Cette fonction construit la chaîne de connexion et teste la connectivité
// avant de retourner l'instance de Connection.
func NewConnection(host string, port int, user, password, dbname string) (*Connection, error) {
	// Construction de la chaîne de connexion PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test de la connexion pour s'assurer qu'elle fonctionne
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Connection{db: db}, nil
}

// Close ferme la connexion à la base de données
func (c *Connection) Close() error {
	return c.db.Close()
}
