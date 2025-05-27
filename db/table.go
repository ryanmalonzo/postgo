package db

import (
	"fmt"
	"reflect"
	"strings"

	"slices"

	_ "github.com/lib/pq"
)

// Create a new table in the database based on the provided model.
// The model must implement the Model interface to provide its table name.
func (c *Connection) CreateTable(model Model) error {
	// Get table name directly from the Model interface
	tableName := model.TableName()

	// Get model metadata using reflection
	metadata, err := GetMetadata(model)
	if err != nil {
		return fmt.Errorf("failed to get model metadata: %w", err)
	}

	// Build column definitions with constraints
	var columns []string
	for _, field := range metadata.Fields {
		colDef := buildColumnDefinition(field)
		columns = append(columns, colDef)
	}

	columnsStr := strings.Join(columns, ", ")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (%s)", tableName, columnsStr)

	_, err = c.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

// buildColumnDefinition creates a SQL column definition string including constraints
func buildColumnDefinition(field Field) string {
	var colType string
	isAutoIncrement := slices.Contains(field.Constraints, "auto_increment")

	// Map Go types to SQL types
	switch field.Type.Kind() {
	case reflect.String:
		colType = "VARCHAR(255)"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isAutoIncrement {
			colType = "SERIAL"
		} else {
			colType = "INTEGER"
		}
	case reflect.Float32, reflect.Float64:
		colType = "FLOAT"
	case reflect.Bool:
		colType = "BOOLEAN"
	default:
		panic(fmt.Sprintf("unsupported field type: %s", field.Type))
	}

	// Escape column name with double quotes
	definition := fmt.Sprintf("\"%s\" %s", field.Name, colType)

	// Add constraints (except auto_increment which is handled via the data type)
	for _, constraint := range field.Constraints {
		switch constraint {
		case "primary_key":
			definition += " PRIMARY KEY"
		case "not_null":
			definition += " NOT NULL"
		case "unique":
			definition += " UNIQUE"
		case "auto_increment":
			// Already handled in the type
		default:
			// For custom constraints, add them directly
			if !strings.Contains(definition, constraint) {
				definition += " " + strings.ToUpper(constraint)
			}
		}
	}

	return definition
}
