package db

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

// CreateTable creates a new table in the database based on the provided schema.
// The schema should be a struct type that defines the table columns.
func (c *Connection) CreateTable(tableName string, schema interface{}) error {
	var columns []string

	// Use reflection to get struct fields and types
	t := reflect.TypeOf(schema)

	// If pointer, get the underlying element
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Ensure we're dealing with a struct
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("schema must be a struct type or pointer to struct, got %v", t.Kind())
	}

	for i := range t.NumField() {
		field := t.Field(i)
		colName := field.Name

		var colType string
		switch field.Type.Kind() {
		case reflect.String:
			colType = "VARCHAR(255)"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			colType = "INT"
		case reflect.Float32, reflect.Float64:
			colType = "FLOAT"
		case reflect.Bool:
			colType = "BOOLEAN"
		default:
			return fmt.Errorf("unsupported column type: %s", field.Type.Name())
		}

		columns = append(columns, fmt.Sprintf("%s %s", colName, colType))
	}

	columnsStr := strings.Join(columns, ", ")

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, columnsStr)
	_, err := c.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
