package query

import (
	"database/sql"
	"strings"
)

type SelectQuery struct {
	BaseQuery
	table   string
	columns []string
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

func (q *SelectQuery) OrderBy(column string, direction string) *SelectQuery {
	q.orderBy = append(q.orderBy, column+" "+direction)
	return q
}

func (q *SelectQuery) Limit(limit int) *SelectQuery {
	q.limit = &limit
	return q
}

func (q *SelectQuery) Offset(offset int) *SelectQuery {
	q.offset = &offset
	return q
}

func (q *SelectQuery) Build() string {
	query := "SELECT " + strings.Join(q.columns, ", ") + " FROM " + q.table
	commonClauses := q.buildCommonClauses()
	if commonClauses != "" {
		query += " " + commonClauses
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
