package query

import (
	"database/sql"
	"fmt"
	"strings"
)

type InsertQuery struct {
	table   string
	columns []string
	values  []interface{} // Changé en interface{} pour supporter tous types
}

func NewInsertQuery(table string) *InsertQuery {
	return &InsertQuery{
		table: table,
	}
}

func (q *InsertQuery) AddColumn(column string) *InsertQuery {
	q.columns = append(q.columns, column)
	return q
}

// AddValue accepte maintenant n'importe quel type
func (q *InsertQuery) AddValue(value interface{}) *InsertQuery {
	q.values = append(q.values, value)
	return q
}

// AddStringValue pour garder la compatibilité avec l'ancienne API
func (q *InsertQuery) AddStringValue(value string) *InsertQuery {
	q.values = append(q.values, value)
	return q
}

func (q *InsertQuery) Build() string {
	if len(q.columns) != len(q.values) {
		panic("Le nombre de colonnes doit être égal au nombre de valeurs")
	}

	columnsList := strings.Join(q.columns, ", ")

	// Créer des placeholders ($1, $2, etc.) pour PostgreSQL
	placeholders := make([]string, len(q.values))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	placeholdersList := strings.Join(placeholders, ", ")

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", q.table, columnsList, placeholdersList)
}

func (q *InsertQuery) Execute(db *sql.DB) (sql.Result, error) {
	query := q.Build()
	fmt.Printf("Exécution de la requête: %s\n", query)
	fmt.Printf("Avec les valeurs: %v\n", q.values)

	result, err := db.Exec(query, q.values...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
