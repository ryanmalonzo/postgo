package query

import (
	"database/sql"
	"fmt"
	"strings"
)

type UpdateQuery struct {
	BaseQuery
	table   string
	columns []string
	values  []interface{}
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

func (q *UpdateQuery) AddValue(value interface{}) *UpdateQuery {
	q.values = append(q.values, value)
	return q
}

func (q *UpdateQuery) Where(condition string) *UpdateQuery {
	q.conditions = append(q.conditions, condition)
	return q
}

func (q *UpdateQuery) OrderBy(column string, direction string) *UpdateQuery {
	q.orderBy = append(q.orderBy, column+" "+direction)
	return q
}

func (q *UpdateQuery) Limit(limit int) *UpdateQuery {
	q.limit = &limit
	return q
}

func (q *UpdateQuery) Offset(offset int) *UpdateQuery {
	q.offset = &offset
	return q
}

func (q *UpdateQuery) Build() string {
	if len(q.columns) != len(q.values) {
		panic("Number of columns and values must match")
	}

	var setPairs []string
	for i, column := range q.columns {
		setPairs = append(setPairs, fmt.Sprintf("%s = $%d", column, i+1))
	}

	query := fmt.Sprintf("UPDATE %s SET %s", q.table, strings.Join(setPairs, ", "))

	commonClauses := q.buildCommonClauses()
	if commonClauses != "" {
		query += " " + commonClauses
	}

	return query
}

func (q *UpdateQuery) Execute(db *sql.DB) error {
	query := q.Build()
	_, err := db.Exec(query, q.values...)
	return err
}

// GetValues retourne les valeurs de la requête (utile pour le générateur)
func (q *UpdateQuery) GetValues() []interface{} {
	return q.values
}
