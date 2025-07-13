package query

import (
	"database/sql"
	"fmt"
	"strings"
)

type SelectQuery struct {
	BaseQuery
	table   string
	columns []string
	values  []interface{}
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
	
	// Gestion manuelle des clauses WHERE avec numérotation des paramètres
	if len(q.conditions) > 0 {
		// Remplacer les $1 par les bons numéros de paramètres
		whereClause := "WHERE " + strings.Join(q.conditions, " AND ")
		query += " " + whereClause
	}
	
	// Ajouter les autres clauses (ORDER BY, LIMIT, OFFSET) qui n'ont pas de paramètres
	if len(q.orderBy) > 0 {
		query += " ORDER BY " + strings.Join(q.orderBy, ", ")
	}
	
	if q.limit != nil {
		query += " LIMIT " + fmt.Sprintf("%d", *q.limit)
	}
	
	if q.offset != nil {
		query += " OFFSET " + fmt.Sprintf("%d", *q.offset)
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
