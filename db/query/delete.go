package query

import (
	"database/sql"
)

type DeleteQuery struct {
	BaseQuery
	table string
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
	commonClauses := q.buildCommonClauses()
	if commonClauses != "" {
		query += " " + commonClauses
	}
	return query
}

func (q *DeleteQuery) Execute(db *sql.DB) error {
	query := q.Build()
	_, err := db.Exec(query)
	return err
}
