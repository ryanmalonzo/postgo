package query

import (
	"database/sql"
	"strings"
)

type SelectQuery struct {
	table      string
	columns    []string
	conditions []string
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

func (q *SelectQuery) Build() string {
	query := "SELECT " + strings.Join(q.columns, ", ") + " FROM " + q.table
	if len(q.conditions) > 0 {
		query += " WHERE " + strings.Join(q.conditions, " AND ")
	}
	return query
}

func (q *SelectQuery) Execute(db *sql.DB) (*sql.Rows, error) {
	query := q.Build()
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
