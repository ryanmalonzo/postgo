package query

import (
	"database/sql"
	"strings"
)

type SelectQuery struct {
	table      string
	columns    []string
	conditions []string
	values     []interface{}
}

func NewSelectQuery(table string) *SelectQuery {
	return &SelectQuery{
		table: table,
	}
}

func (q *SelectQuery) AddColumn(column string) *SelectQuery {
	q.columns = append(q.columns, column)
	return q
}

func (q *SelectQuery) AddCondition(condition string) *SelectQuery {
	q.conditions = append(q.conditions, condition)
	return q
}

func (q *SelectQuery) Where(condition string) *SelectQuery {
	q.conditions = append(q.conditions, condition)
	return q
}

func (q *SelectQuery) Build() string {
	query := "SELECT " + strings.Join(q.columns, ", ") + " FROM " + q.table
	
	// Gestion des clauses WHERE
	if len(q.conditions) > 0 {
		whereClause := "WHERE " + strings.Join(q.conditions, " AND ")
		query += " " + whereClause
	}
	
	return query
}

func (q *SelectQuery) Execute(db *sql.DB) (*sql.Rows, error) {
	query := q.Build()
	rows, err := db.Query(query, q.values...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// WhereWithValue ajoute une condition WHERE avec une valeur paramétrée
func (q *SelectQuery) WhereWithValue(condition string, value interface{}) *SelectQuery {
	q.conditions = append(q.conditions, condition)
	q.values = append(q.values, value)
	return q
}

// GetValues retourne les valeurs de la requête (utile pour le générateur)
func (q *SelectQuery) GetValues() []interface{} {
	return q.values
}
