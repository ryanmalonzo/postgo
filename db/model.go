package db

import (
	"fmt"
	"reflect"
	"strings"
)

// Base interface to be implemented by all models
type Model interface {
	TableName() string
}

type BaseModel struct {
	ID int64 `db:"primary_key,auto_increment"`
}

type Field struct {
	Name        string
	Type        reflect.Type
	Constraints []string
}

type ModelMetadata struct {
	TableName string
	Fields    []Field
}

func GetMetadata(model any) (*ModelMetadata, error) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct type or pointer to struct, got %v", t.Kind())
	}

	var fields []Field
	if err := extractFields(t, &fields); err != nil {
		return nil, err
	}

	return &ModelMetadata{
		TableName: strings.ToLower(t.Name()),
		Fields:    fields,
	}, nil
}

// tag = `db:"primary_key,auto_increment"` for example
func parseConstraints(tag string) []string {
	constraints := []string{}
	if tag == "" {
		return constraints
	}

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			constraints = append(constraints, part)
		}
	}

	return constraints
}

// Recursively extract and store fields from a struct type
func extractFields(t reflect.Type, fields *[]Field) error {
	for i := range t.NumField() {
		field := t.Field(i)
		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// Handle embedded structs recursively (such as BaseModel)
		if field.Anonymous && fieldType.Kind() == reflect.Struct {
			if err := extractFields(fieldType, fields); err != nil {
				return err
			}
			continue
		}

		constraints := parseConstraints(field.Tag.Get("db"))

		*fields = append(*fields, Field{
			Name:        strings.ToLower(field.Name),
			Type:        fieldType,
			Constraints: constraints,
		})
	}

	return nil
}
