package query

import (
	"fmt"
	"strings"
)

// QueryBuilder interface commune pour toutes les requêtes
type QueryBuilder interface {
	Build() string
}

// BaseQuery contient les champs communs à toutes les requêtes
type BaseQuery struct {
	conditions []string
	orderBy    []string
	limit      *int
	offset     *int
}

// Where ajoute une condition WHERE à la requête
func (q *BaseQuery) Where(condition string) *BaseQuery {
	q.conditions = append(q.conditions, condition)
	return q
}

// OrderBy ajoute une clause ORDER BY
func (q *BaseQuery) OrderBy(column string, direction string) *BaseQuery {
	q.orderBy = append(q.orderBy, column+" "+direction)
	return q
}

// Limit définit la limite de résultats
func (q *BaseQuery) Limit(limit int) *BaseQuery {
	q.limit = &limit
	return q
}

// Offset définit le décalage des résultats
func (q *BaseQuery) Offset(offset int) *BaseQuery {
	q.offset = &offset
	return q
}

// buildCommonClauses construit les clauses communes (WHERE, ORDER BY, LIMIT, OFFSET)
func (q *BaseQuery) buildCommonClauses() string {
	var clauses []string

	if len(q.conditions) > 0 {
		clauses = append(clauses, "WHERE "+strings.Join(q.conditions, " AND "))
	}

	if len(q.orderBy) > 0 {
		clauses = append(clauses, "ORDER BY "+strings.Join(q.orderBy, ", "))
	}

	if q.limit != nil {
		clauses = append(clauses, "LIMIT "+fmt.Sprintf("%d", *q.limit))
	}

	if q.offset != nil {
		clauses = append(clauses, "OFFSET "+fmt.Sprintf("%d", *q.offset))
	}

	return strings.Join(clauses, " ")
}
