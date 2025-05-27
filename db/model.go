package db

import (
	"fmt"
	"reflect"
	"strings"
)

// Interface de base à implémenter par tous les modèles
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

// GetMetadata extrait les métadonnées d'un modèle en utilisant la réflexion
// pour analyser sa structure et ses tags. Cette fonction examine la structure
// d'un type Go et retourne les informations nécessaires pour créer une table SQL.
func GetMetadata(model any) (*ModelMetadata, error) {
	t := reflect.TypeOf(model)
	// Si c'est un pointeur, on récupère le type sous-jacent
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct type or pointer to struct, got %v", t.Kind())
	}

	var fields []Field
	// Extraction récursive des champs du struct
	if err := extractFields(t, &fields); err != nil {
		return nil, err
	}

	return &ModelMetadata{
		TableName: strings.ToLower(t.Name()),
		Fields:    fields,
	}, nil
}

// parseConstraints analyse les tags `db` pour extraire les contraintes SQL
// Exemple: `db:"primary_key,auto_increment"` devient ["primary_key", "auto_increment"]
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

// extractFields parcourt récursivement les champs d'un struct en utilisant la réflexion.
// Cette fonction gère les structs imbriqués (comme BaseModel) et extrait les informations
// de type et les contraintes de chaque champ pour construire le schéma de la table.
func extractFields(t reflect.Type, fields *[]Field) error {
	// Parcours de tous les champs du struct
	for i := range t.NumField() {
		field := t.Field(i)
		fieldType := field.Type

		// Déréférencement des pointeurs pour obtenir le type réel
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// Gestion des structs imbriqués (embedded structs) comme BaseModel
		// Les champs anonymes sont traités récursivement
		if field.Anonymous && fieldType.Kind() == reflect.Struct {
			if err := extractFields(fieldType, fields); err != nil {
				return err
			}
			continue
		}

		// Extraction des contraintes depuis le tag `db`
		constraints := parseConstraints(field.Tag.Get("db"))

		*fields = append(*fields, Field{
			Name:        strings.ToLower(field.Name),
			Type:        fieldType,
			Constraints: constraints,
		})
	}

	return nil
}
