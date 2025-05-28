package query

import (
	"database/sql"
	"fmt"
	"strings"
)

type UpdateQuery struct {
	table      string
	columns    []string
	conditions []string
	values     []interface{}
}

func NewUpdateQuery(table string) *UpdateQuery {
	return &UpdateQuery{
		table: table,
	}
}

func (q *UpdateQuery) AddColumn(column string) *UpdateQuery {
	q.columns = append(q.columns, column)
	return q
}

func (q *UpdateQuery) AddCondition(condition string) *UpdateQuery {
	q.conditions = append(q.conditions, condition)
	return q
}

func (q *UpdateQuery) AddValue(value interface{}) *UpdateQuery {
	q.values = append(q.values, value)
	return q
}

func (q *UpdateQuery) Build() string {
	// Construire les paires colonne=valeur
	setClauses := make([]string, len(q.columns))
	for i, col := range q.columns {
		setClauses[i] = col + " = $" + fmt.Sprintf("%d", i+1)
	}

	query := "UPDATE " + q.table + " SET " + strings.Join(setClauses, ", ")
	if len(q.conditions) > 0 {
		query += " WHERE " + strings.Join(q.conditions, " AND ")
	}
	return query
}

func (q *UpdateQuery) Execute(db *sql.DB) error {
	query := q.Build()
	_, err := db.Exec(query, q.values...)
	return err
}
