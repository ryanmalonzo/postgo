package functions

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func CreateTable(db *sql.DB, tableName string, schema interface{}) error {
	var columns []string
	
	// Use reflection to get struct fields and types
	t := reflect.TypeOf(schema)
	
	// If pointer, get the element type
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
	
	// Join the column definitions into a comma-separated string
	columnsStr := strings.Join(columns, ", ")

	// Create a table with the specified name
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, columnsStr)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
