package db

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

// AttributeType représente les types de données supportés
type AttributeType string

const (
	String  AttributeType = "VARCHAR(255)"
	Integer AttributeType = "INTEGER"
	Float   AttributeType = "FLOAT"
	Boolean AttributeType = "BOOLEAN"
)

// Attribute représente une colonne de table avec ses contraintes
type Attribute struct {
	name        string
	dataType    AttributeType
	constraints []string
}

// AttributeBuilder permet de construire un attribut avec le pattern builder
type AttributeBuilder struct {
	attribute   *Attribute
	tableBuilder *TableBuilder // Référence vers le table builder parent
}

// NotNull ajoute la contrainte NOT NULL à l'attribut
func (ab *AttributeBuilder) NotNull() *AttributeBuilder {
	ab.attribute.constraints = append(ab.attribute.constraints, "NOT NULL")
	return ab
}

// Unique ajoute la contrainte UNIQUE à l'attribut
func (ab *AttributeBuilder) Unique() *AttributeBuilder {
	ab.attribute.constraints = append(ab.attribute.constraints, "UNIQUE")
	return ab
}

// Build finalise la construction de l'attribut et l'ajoute à la table
func (ab *AttributeBuilder) Build() *TableBuilder {
	if ab.tableBuilder != nil {
		ab.tableBuilder.attributes = append(ab.tableBuilder.attributes, ab.attribute)
		return ab.tableBuilder
	}
	return nil
}

// TableBuilder permet de construire une table avec le pattern builder
type TableBuilder struct {
	name       string
	attributes []*Attribute
}

// NewTable crée un nouveau builder de table avec l'ID auto-incrémenté obligatoire
func NewTable(name string) *TableBuilder {
	tb := &TableBuilder{
		name:       name,
		attributes: make([]*Attribute, 0),
	}
	
	// Ajout automatique de l'ID auto-incrémenté (obligatoire)
	idAttribute := &Attribute{
		name:        "id",
		dataType:    "SERIAL",
		constraints: []string{"PRIMARY KEY"},
	}
	tb.attributes = append(tb.attributes, idAttribute)
	
	return tb
}

// AddAttribute ajoute un attribut à la table avec une syntaxe fluide
func (tb *TableBuilder) AddAttribute(name string, dataType AttributeType) *AttributeBuilder {
	attr := &Attribute{
		name:        name,
		dataType:    dataType,
		constraints: make([]string, 0),
	}
	
	return &AttributeBuilder{
		attribute:    attr,
		tableBuilder: tb,
	}
}

// BuildSQL finalise la construction de la table et retourne la requête SQL
func (tb *TableBuilder) BuildSQL() string {
	var columns []string
	
	for _, attr := range tb.attributes {
		definition := fmt.Sprintf("\"%s\" %s", attr.name, attr.dataType)
		
		// Ajout des contraintes
		for _, constraint := range attr.constraints {
			definition += " " + constraint
		}
		
		columns = append(columns, definition)
	}
	
	columnsStr := strings.Join(columns, ", ")
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (%s)", tb.name, columnsStr)
}

// CreateTable crée une nouvelle table dans la base de données en utilisant un TableBuilder
func (c *Connection) CreateTable(tableBuilder *TableBuilder) error {
	query := tableBuilder.BuildSQL()
	
	_, err := c.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

// === MÉTHODES POUR LE GÉNÉRATEUR ===

// GetName retourne le nom de la table
func (tb *TableBuilder) GetName() string {
	return tb.name
}

// GetAttributes retourne tous les attributs de la table
func (tb *TableBuilder) GetAttributes() []*Attribute {
	return tb.attributes
}

// GetName retourne le nom de l'attribut
func (a *Attribute) GetName() string {
	return a.name
}

// GetDataType retourne le type de données de l'attribut
func (a *Attribute) GetDataType() AttributeType {
	return a.dataType
}

// GetConstraints retourne les contraintes de l'attribut
func (a *Attribute) GetConstraints() []string {
	return a.constraints
}

// IsRequired vérifie si l'attribut a la contrainte NOT NULL
func (a *Attribute) IsRequired() bool {
	for _, constraint := range a.constraints {
		if constraint == "NOT NULL" {
			return true
		}
	}
	return false
}

// IsUnique vérifie si l'attribut a la contrainte UNIQUE
func (a *Attribute) IsUnique() bool {
	for _, constraint := range a.constraints {
		if constraint == "UNIQUE" {
			return true
		}
	}
	return false
}

// GetGoType retourne le type Go correspondant au type de données
func (a *Attribute) GetGoType() string {
	switch a.dataType {
	case String:
		return "string"
	case Integer:
		return "int"
	case Float:
		return "float64"
	case Boolean:
		return "bool"
	case "SERIAL":
		return "int"
	default:
		return "interface{}"
	}
}
