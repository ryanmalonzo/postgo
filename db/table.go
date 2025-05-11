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

	// Process all fields, including embedded structs and their fields, if any (such as BaseModel)
	if err := processFields(t, "", &columns); err != nil {
		return err
	}

	columnsStr := strings.Join(columns, ", ")

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, columnsStr)
	_, err := c.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func processFields(t reflect.Type, prefix string, columns *[]string) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// If field is an embedded struct, process its fields recursively
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if err := processFields(field.Type, prefix, columns); err != nil {
				return err
			}
			continue
		}

		colName := field.Name
		if prefix != "" {
			colName = prefix + colName
		}

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

		*columns = append(*columns, fmt.Sprintf("%s %s", colName, colType))
	}
	return nil
}
