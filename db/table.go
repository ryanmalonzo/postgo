package db

import (
	"fmt"
	"reflect"
	"strings"

	"slices"

	_ "github.com/lib/pq"
)

// CreateTable crée une nouvelle table dans la base de données basée sur le modèle fourni.
// Le modèle doit implémenter l'interface Model pour fournir le nom de sa table.
// Cette fonction utilise la réflexion pour analyser la structure du modèle et
// générer automatiquement le schéma SQL correspondant.
func (c *Connection) CreateTable(model Model) error {
	// Récupération du nom de table depuis l'interface Model
	tableName := model.TableName()

	// Extraction des métadonnées du modèle via réflexion
	metadata, err := GetMetadata(model)
	if err != nil {
		return fmt.Errorf("failed to get model metadata: %w", err)
	}

	// Construction des définitions de colonnes avec leurs contraintes
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

// buildColumnDefinition construit une définition de colonne SQL complète avec ses contraintes.
// Cette fonction mappe les types Go vers les types SQL PostgreSQL et applique
// les contraintes définies dans les tags du struct.
func buildColumnDefinition(field Field) string {
	var colType string
	isAutoIncrement := slices.Contains(field.Constraints, "auto_increment")

	// Mapping des types Go vers les types SQL PostgreSQL
	switch field.Type.Kind() {
	case reflect.String:
		colType = "VARCHAR(255)"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isAutoIncrement {
			colType = "SERIAL" // Type PostgreSQL pour auto-incrémentation
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

	// Échappement du nom de colonne avec des guillemets doubles
	definition := fmt.Sprintf("\"%s\" %s", field.Name, colType)

	// Ajout des contraintes (sauf auto_increment qui est géré par le type de données)
	for _, constraint := range field.Constraints {
		switch constraint {
		case "primary_key":
			definition += " PRIMARY KEY"
		case "not_null":
			definition += " NOT NULL"
		case "unique":
			definition += " UNIQUE"
		case "auto_increment":
			// Déjà géré dans le type de données (SERIAL)
		default:
			// Pour les contraintes personnalisées, les ajouter directement
			if !strings.Contains(definition, constraint) {
				definition += " " + strings.ToUpper(constraint)
			}
		}
	}

	return definition
}
