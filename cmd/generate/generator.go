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
	
	// Générer les composants pour Select
	selectMethods := generateSelectComponents(attributes, titleName, tableName)
	
	// Générer le struct principal
	mainStruct := generateMainStruct(attributes, titleName, tableName)
	
	content := fmt.Sprintf(`// Code généré automatiquement - NE PAS MODIFIER
package generated

import (
	"database/sql"
	"fmt"
	"postgo/db"
	"postgo/db/query"
)

%s// %sTable représente la table %s
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

// %sSelectBuilder permet de sélectionner des données de la table %s
type %sSelectBuilder struct {
	query *query.SelectQuery
}

// %sSelectResult représente les résultats possibles d'une sélection
type %sSelectResult struct {
	selectedColumns []string
	query           *query.SelectQuery
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

// Select crée un nouveau builder pour sélectionner dans la table %s
func (t *%sTable) Select() *%sSelectBuilder {
	return &%sSelectBuilder{
		query: query.NewSelectQuery("%s"),
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
%s
`,
		mainStruct,            // main struct
		titleName, tableName,  // Table comment
		titleName,             // type Table struct
		tableName,             // instance comment
		titleName, titleName, tableName,  // var Table = &Table{Name:}
		titleName, tableName,  // InsertBuilder comment
		titleName,             // type InsertBuilder struct
		insertFields,          // insert fields
		titleName, tableName,  // UpdateBuilder comment
		titleName,             // type UpdateBuilder struct
		updateFields,          // update fields
		titleName, tableName,  // DeleteBuilder comment
		titleName,             // type DeleteBuilder struct
		titleName, tableName,  // SelectBuilder comment
		titleName,             // type SelectBuilder struct
		titleName,             // SelectResult comment
		titleName,             // type SelectResult struct
		tableName,             // Insert() comment
		titleName, titleName,  // Insert() function
		titleName, tableName,  // return &InsertBuilder{NewInsertQuery()}
		tableName,             // Update() comment
		titleName, titleName,  // Update() function
		titleName, tableName,  // return &UpdateBuilder{NewUpdateQuery()}
		tableName,             // Delete() comment
		titleName, titleName,  // Delete() function
		titleName, tableName,  // return &DeleteBuilder{NewDeleteQuery()}
		tableName,             // Select() comment
		titleName, titleName,  // Select() function
		titleName, tableName,  // return &SelectBuilder{NewSelectQuery()}
		insertMethods,         // insert methods
		updateMethods,         // update methods
		titleName,             // Execute() comment for insert
		insertRequiredChecks,  // required checks for insert
		titleName,             // Build() for insert
		titleName, titleName,  // Where() for update
		titleName,             // Execute() for update
		titleName,             // Build() for update
		titleName, titleName,  // Where() for delete
		titleName,             // Execute() for delete
		titleName,             // Build() for delete
		selectMethods)         // select methods

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

// generateSelectComponents génère les composants pour le builder de sélection
func generateSelectComponents(attributes []*db.Attribute, titleName, tableName string) string {
	var selectMethods []string
	
	// Utiliser le nom singulier pour le struct (ex: User au lieu de Users)
	singularName := titleName
	if strings.HasSuffix(titleName, "s") {
		singularName = titleName[:len(titleName)-1]
	}
	
	// Importer database/sql est nécessaire pour sql.ErrNoRows
	
	// Méthode SelectAll
	selectMethods = append(selectMethods, fmt.Sprintf(`
// SelectAll sélectionne toutes les colonnes de la table %s
func (b *%sSelectBuilder) SelectAll() *%sSelectResult {
	b.query.AddColumn("*")
	return &%sSelectResult{
		selectedColumns: []string{"*"},
		query:           b.query,
	}
}`, titleName, titleName, titleName, titleName))

	// Méthodes pour chaque colonne sur le SelectBuilder
	for _, attr := range attributes {
		attrName := attr.GetName()
		titleAttrName := toCamelCase(attrName)
		
		selectMethods = append(selectMethods, fmt.Sprintf(`
// Select%s sélectionne la colonne %s
func (b *%sSelectBuilder) Select%s() *%sSelectResult {
	b.query.AddColumn("%s")
	return &%sSelectResult{
		selectedColumns: []string{"%s"},
		query:           b.query,
	}
}`, titleAttrName, attrName, titleName, titleAttrName, titleName, attrName, titleName, attrName))
	}

	// Méthodes pour chaque colonne sur le SelectResult (pour permettre l'enchaînement)
	for _, attr := range attributes {
		attrName := attr.GetName()
		titleAttrName := toCamelCase(attrName)
		
		selectMethods = append(selectMethods, fmt.Sprintf(`
// Select%s ajoute la colonne %s à la sélection
func (r *%sSelectResult) Select%s() *%sSelectResult {
	r.query.AddColumn("%s")
	r.selectedColumns = append(r.selectedColumns, "%s")
	return r
}`, titleAttrName, attrName, titleName, titleAttrName, titleName, attrName, attrName))
	}

	// Méthode SelectColumns générique
	selectMethods = append(selectMethods, fmt.Sprintf(`
// SelectColumns sélectionne plusieurs colonnes spécifiques
func (b *%sSelectBuilder) SelectColumns(columns ...string) *%sSelectResult {
	for _, column := range columns {
		b.query.AddColumn(column)
	}
	return &%sSelectResult{
		selectedColumns: columns,
		query:           b.query,
	}
}`, titleName, titleName, titleName))

	// Méthodes WHERE typées pour SelectResult
	var whereMethods []string
	
	whereMethods = append(whereMethods, fmt.Sprintf(`
// Where ajoute une condition WHERE à la requête de sélection
func (r *%sSelectResult) Where(condition string) *%sSelectResult {
	r.query.Where(condition)
	return r
}`, titleName, titleName))

	for _, attr := range attributes {
		attrName := attr.GetName()
		titleAttrName := toCamelCase(attrName)
		goType := attr.GetGoType()
		
		whereMethods = append(whereMethods, fmt.Sprintf(`
// Where%s ajoute une condition WHERE sur %s
func (r *%sSelectResult) Where%s(%s %s) *%sSelectResult {
	r.query.WhereWithValue("%s = $1", %s)
	return r
}`, titleAttrName, attrName, titleName, titleAttrName, strings.ToLower(attrName), goType, titleName, attrName, strings.ToLower(attrName)))
	}

	// Méthodes Execute
	executeMethods := fmt.Sprintf(`
// Execute exécute la requête et retourne les résultats typés
func (r *%sSelectResult) Execute(conn *db.Connection) ([]%s, error) {
	sqlQuery := r.query.Build()
	values := r.query.GetValues()
	
	rows, err := conn.GetDB().Query(sqlQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var results []%s
	
	for rows.Next() {
		var result %s
		
		// Si on sélectionne toutes les colonnes ou certaines colonnes spécifiques
		if len(r.selectedColumns) == 1 && r.selectedColumns[0] == "*" {
			// Scan toutes les colonnes
			err = rows.Scan(%s)
		} else {
			// Scan seulement les colonnes sélectionnées
			var scanTargets []interface{}
			for _, col := range r.selectedColumns {
				switch col {
%s
				}
			}
			err = rows.Scan(scanTargets...)
		}
		
		if err != nil {
			return nil, err
		}
		
		results = append(results, result)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return results, nil
}

// ExecuteOne exécute la requête et retourne un seul résultat
func (r *%sSelectResult) ExecuteOne(conn *db.Connection) (*%s, error) {
	results, err := r.Execute(conn)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}
	
	return &results[0], nil
}

// Build retourne la requête SQL pour la sélection
func (r *%sSelectResult) Build() (string, []interface{}) {
	return r.query.Build(), r.query.GetValues()
}`, titleName, singularName, singularName, singularName, generateAllColumnsScan(attributes), generateColumnCases(attributes, singularName), titleName, singularName, titleName)

	return strings.Join(selectMethods, "") + strings.Join(whereMethods, "") + executeMethods
}

// generateAllColumnsScan génère le code pour scanner toutes les colonnes
func generateAllColumnsScan(attributes []*db.Attribute) string {
	var scans []string
	for _, attr := range attributes {
		attrName := attr.GetName()
		scans = append(scans, "&result."+toCamelCase(attrName))
	}
	return strings.Join(scans, ", ")
}

// generateColumnCases génère les cases pour le switch des colonnes
func generateColumnCases(attributes []*db.Attribute, structName string) string {
	var cases []string
	for _, attr := range attributes {
		attrName := attr.GetName()
		titleAttrName := toCamelCase(attrName)
		cases = append(cases, fmt.Sprintf(`				case "%s":
					scanTargets = append(scanTargets, &result.%s)`, attrName, titleAttrName))
	}
	return strings.Join(cases, "\n")
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

// generateMainStruct génère le struct principal représentant une ligne de la table
func generateMainStruct(attributes []*db.Attribute, titleName, tableName string) string {
	var fields []string
	
	// Utiliser le nom singulier pour le struct (ex: User au lieu de Users)
	singularName := titleName
	if strings.HasSuffix(titleName, "s") {
		singularName = titleName[:len(titleName)-1]
	}
	
	for _, attr := range attributes {
		attrName := attr.GetName()
		titleAttrName := toCamelCase(attrName)
		goType := attr.GetGoType()
		
		fields = append(fields, fmt.Sprintf("	%s %s `db:\"%s\"`", titleAttrName, goType, attrName))
	}
	
	return fmt.Sprintf(`// %s représente une ligne de la table %s
type %s struct {
%s
}

`, singularName, tableName, singularName, strings.Join(fields, "\n"))
}
