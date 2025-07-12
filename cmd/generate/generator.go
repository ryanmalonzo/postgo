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
	// Obtenir les attributs de la table
	attributes := table.GetAttributes()
	titleName := titleCase(tableName)
	
	var setMethods []string
	var requiredChecks []string
	var fieldDeclarations []string
	
	for _, attr := range attributes {
		attrName := attr.GetName()
		
		// Ignorer l'ID car il est auto-généré
		if attrName == "id" {
			continue
		}
		
		// Convertir le nom de l'attribut en format Go (snake_case vers CamelCase)
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
			requiredChecks = append(requiredChecks, requiredCheck)
		}
	}
	
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

// Insert crée un nouveau builder pour insérer dans la table %s
func (t *%sTable) Insert() *%sInsertBuilder {
	return &%sInsertBuilder{
		query: query.NewInsertQuery("%s"),
	}
}
%s

// Execute exécute la requête d'insertion
func (b *%sInsertBuilder) Execute(conn *db.Connection) error {
%s
	sqlQuery, args := b.query.Build(), b.query.GetValues()
	_, err := conn.GetDB().Exec(sqlQuery, args...)
	return err
}

// Build retourne la requête SQL et les arguments
func (b *%sInsertBuilder) Build() (string, []interface{}) {
	return b.query.Build(), b.query.GetValues()
}
`,
		titleName, tableName,
		titleName,
		tableName,
		titleName, titleName, tableName,
		titleName, tableName,
		titleName,
		strings.Join(fieldDeclarations, "\n"),
		tableName,
		titleName, titleName,
		titleName, tableName,
		strings.Join(setMethods, ""),
		titleName,
		strings.Join(requiredChecks, "\n"),
		titleName)

	return writeFile(filepath.Join(outputDir, tableName+".go"), content)
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
