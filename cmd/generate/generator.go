package main

import (
	"fmt"
	"os"
	"path/filepath"
	"postgo/db"
	"strings"
)

// generateMainTypes génère le fichier avec les types principaux
func generateMainTypes(outputDir string, tables map[string]*db.TableBuilder) error {
	content := `// Code généré automatiquement - NE PAS MODIFIER
package generated

import (
	"postgo/db"
)

// Types de base pour la validation
type TableReference struct {
	Name string
}

// Interface commune pour tous les builders d'insertion
type InsertBuilder interface {
	Execute(conn *db.Connection) error
	Build() (string, []interface{})
}

// Interface commune pour tous les builders d'update
type UpdateBuilder interface {
	Execute(conn *db.Connection) error
	Build() (string, []interface{})
	Where(condition string) UpdateBuilder
}

// Interface commune pour tous les builders de suppression
type DeleteBuilder interface {
	Execute(conn *db.Connection) error
	Build() (string, []interface{})
	Where(condition string) DeleteBuilder
}
`

	return writeFile(filepath.Join(outputDir, "types.go"), content)
}

// generateIndexFile génère le fichier d'index qui exporte toutes les tables
func generateIndexFile(outputDir string, tables map[string]*db.TableBuilder) error {
	var exports []string

	for tableName := range tables {
		exports = append(exports, fmt.Sprintf(`	%s = %s.Table`, titleCase(tableName), titleCase(tableName)))
	}

	content := fmt.Sprintf(`// Code généré automatiquement - NE PAS MODIFIER
package generated

var (
%s
)
`, strings.Join(exports, "\n"))

	return writeFile(filepath.Join(outputDir, "tables.go"), content)
}

// generateTableFile génère un fichier dédié pour chaque table
func generateTableFile(outputDir, tableName string, table *db.TableBuilder) error {
	attributes := table.GetAttributes()
	titleName := titleCase(tableName)
	
	// Générer les composants pour Insert
	insertFields, insertMethods, insertRequiredChecks := generateInsertComponents(attributes, titleName)
	
	// Générer les composants pour Update
	updateFields, updateMethods := generateUpdateComponents(attributes, titleName)
	
	content := fmt.Sprintf(`// Code généré automatiquement - NE PAS MODIFIER
package generated

import (
	"fmt"
	"postgo/db"
	"postgo/db/query"
)

// %sTable représente la table %s
type %sTable struct {
	Name string
}

// Instance globale de la table %s
var %s = &%sTable{
	Name: "%s",
}

// %sInsertBuilder permet d'insérer des données dans la table %s
type %sInsertBuilder struct {
	query *query.InsertQuery
%s
}

// %sUpdateBuilder permet de mettre à jour des données dans la table %s
type %sUpdateBuilder struct {
	query *query.UpdateQuery
%s
}

// %sDeleteBuilder permet de supprimer des données de la table %s
type %sDeleteBuilder struct {
	query *query.DeleteQuery
}

// Insert crée un nouveau builder pour insérer dans la table %s
func (t *%sTable) Insert() *%sInsertBuilder {
	return &%sInsertBuilder{
		query: query.NewInsertQuery("%s"),
	}
}

// Update crée un nouveau builder pour mettre à jour la table %s
func (t *%sTable) Update() *%sUpdateBuilder {
	return &%sUpdateBuilder{
		query: query.NewUpdateQuery("%s"),
	}
}

// Delete crée un nouveau builder pour supprimer de la table %s
func (t *%sTable) Delete() *%sDeleteBuilder {
	return &%sDeleteBuilder{
		query: query.NewDeleteQuery("%s"),
	}
}
%s
%s

// Execute exécute la requête d'insertion
func (b *%sInsertBuilder) Execute(conn *db.Connection) error {
%s
	sqlQuery, args := b.query.Build(), b.query.GetValues()
	_, err := conn.GetDB().Exec(sqlQuery, args...)
	return err
}

// Build retourne la requête SQL et les arguments pour l'insertion
func (b *%sInsertBuilder) Build() (string, []interface{}) {
	return b.query.Build(), b.query.GetValues()
}

// Where ajoute une condition WHERE à la requête d'update
func (b *%sUpdateBuilder) Where(condition string) *%sUpdateBuilder {
	b.query.Where(condition)
	return b
}

// Execute exécute la requête d'update
func (b *%sUpdateBuilder) Execute(conn *db.Connection) error {
	if len(b.query.GetValues()) == 0 {
		return fmt.Errorf("aucune colonne à mettre à jour")
	}
	sqlQuery, args := b.query.Build(), b.query.GetValues()
	_, err := conn.GetDB().Exec(sqlQuery, args...)
	return err
}

// Build retourne la requête SQL et les arguments pour l'update
func (b *%sUpdateBuilder) Build() (string, []interface{}) {
	return b.query.Build(), b.query.GetValues()
}

// Where ajoute une condition WHERE à la requête de suppression
func (b *%sDeleteBuilder) Where(condition string) *%sDeleteBuilder {
	b.query.AddCondition(condition)
	return b
}

// Execute exécute la requête de suppression
func (b *%sDeleteBuilder) Execute(conn *db.Connection) error {
	sqlQuery := b.query.Build()
	_, err := conn.GetDB().Exec(sqlQuery)
	return err
}

// Build retourne la requête SQL pour la suppression
func (b *%sDeleteBuilder) Build() (string, []interface{}) {
	return b.query.Build(), []interface{}{}
}
`,
		titleName, tableName,
		titleName,
		tableName,
		titleName, titleName, tableName,
		titleName, tableName,
		titleName,
		insertFields,
		titleName, tableName,
		titleName,
		updateFields,
		titleName, tableName,
		titleName,
		tableName,
		titleName, titleName,
		titleName, tableName,
		tableName,
		titleName, titleName,
		titleName, tableName,
		tableName,
		titleName, titleName,
		titleName, tableName,
		insertMethods,
		updateMethods,
		titleName,
		insertRequiredChecks,
		titleName,
		titleName, titleName,
		titleName,
		titleName,
		titleName, titleName,
		titleName,
		titleName)

	return writeFile(filepath.Join(outputDir, tableName+".go"), content)
}

// generateInsertComponents génère les composants pour le builder d'insertion
func generateInsertComponents(attributes []*db.Attribute, titleName string) (fields, methods, requiredChecks string) {
	var fieldDeclarations []string
	var setMethods []string
	var requiredChecksList []string
	
	for _, attr := range attributes {
		attrName := attr.GetName()
		
		// Ignorer l'ID car il est auto-généré pour les inserts
		if attrName == "id" {
			continue
		}
		
		titleAttrName := toCamelCase(attrName)
		lowerAttrName := strings.ToLower(strings.ReplaceAll(attrName, "_", ""))
		goType := attr.GetGoType()
		
		// Déclaration du champ pour suivre si la valeur a été définie
		fieldDeclarations = append(fieldDeclarations, fmt.Sprintf("	%sSet bool", lowerAttrName))
		
		// Méthode Set pour cet attribut
		setMethod := fmt.Sprintf(`
// Set%s définit la valeur pour la colonne %s
func (b *%sInsertBuilder) Set%s(value %s) *%sInsertBuilder {
	if b.%sSet {
		panic("La colonne %s a déjà été définie")
	}
	b.query.AddColumn("%s").AddValue(value)
	b.%sSet = true
	return b
}`, titleAttrName, attrName, titleName, titleAttrName, goType, titleName, lowerAttrName, attrName, attrName, lowerAttrName)
		
		setMethods = append(setMethods, setMethod)
		
		// Vérification pour les champs obligatoires
		if attr.IsRequired() {
			requiredCheck := fmt.Sprintf(`	if !b.%sSet {
		return fmt.Errorf("la colonne obligatoire '%s' n'a pas été définie")
	}`, lowerAttrName, attrName)
			requiredChecksList = append(requiredChecksList, requiredCheck)
		}
	}
	
	return strings.Join(fieldDeclarations, "\n"), 
		   strings.Join(setMethods, ""), 
		   strings.Join(requiredChecksList, "\n")
}

// generateUpdateComponents génère les composants pour le builder d'update
func generateUpdateComponents(attributes []*db.Attribute, titleName string) (fields, methods string) {
	var fieldDeclarations []string
	var setMethods []string
	
	for _, attr := range attributes {
		attrName := attr.GetName()
		
		// Ignorer l'ID car on ne devrait pas l'updater
		if attrName == "id" {
			continue
		}
		
		titleAttrName := toCamelCase(attrName)
		lowerAttrName := strings.ToLower(strings.ReplaceAll(attrName, "_", ""))
		goType := attr.GetGoType()
		
		// Déclaration du champ pour suivre si la valeur a été définie
		fieldDeclarations = append(fieldDeclarations, fmt.Sprintf("	%sSet bool", lowerAttrName))
		
		// Méthode Set pour cet attribut (pour Update)
		setMethod := fmt.Sprintf(`
// Set%s définit la valeur pour la colonne %s dans l'update
func (b *%sUpdateBuilder) Set%s(value %s) *%sUpdateBuilder {
	if b.%sSet {
		panic("La colonne %s a déjà été définie")
	}
	b.query.AddColumn("%s").AddValue(value)
	b.%sSet = true
	return b
}`, titleAttrName, attrName, titleName, titleAttrName, goType, titleName, lowerAttrName, attrName, attrName, lowerAttrName)
		
		setMethods = append(setMethods, setMethod)
	}
	
	return strings.Join(fieldDeclarations, "\n"), strings.Join(setMethods, "")
}

// toCamelCase convertit une chaîne snake_case en CamelCase
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = capitalize(part)
	}
	return strings.Join(parts, "")
}

// capitalize met en majuscule la première lettre d'une chaîne
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// titleCase remplace strings.Title deprecated
func titleCase(s string) string {
	return capitalize(s)
}

// writeFile écrit le contenu dans un fichier
func writeFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
