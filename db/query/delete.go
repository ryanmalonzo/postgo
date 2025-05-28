package query

import (
	"database/sql"
	"strings"
)

type DeleteQuery struct {
	table      string
	conditions []string
}

func NewDeleteQuery(table string) *DeleteQuery {
	return &DeleteQuery{
		table: table,
	}
}

func (q *DeleteQuery) AddCondition(condition string) *DeleteQuery {
	q.conditions = append(q.conditions, condition)
	return q
}

func (q *DeleteQuery) Build() string {
	query := "DELETE FROM " + q.table
	if len(q.conditions) > 0 {
		query += " WHERE " + strings.Join(q.conditions, " AND ")
	}
	return query
}

func (q *DeleteQuery) Execute(db *sql.DB) error {
	query := q.Build()
	_, err := db.Exec(query)
	return err
}
