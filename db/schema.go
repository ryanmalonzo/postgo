package db

import (
	"fmt"
	"postgo/logging"
)

// Schema représente le registre global des tables
type Schema struct {
	tables map[string]*TableBuilder
	order  []string // Pour maintenir l'ordre de création
}

// Instance globale du schéma
var globalSchema *Schema

// init initialise automatiquement le schéma avec toutes les tables
func init() {
	globalSchema = &Schema{
		tables: make(map[string]*TableBuilder),
		order:  make([]string, 0),
	}

	// Enregistrement automatique de toutes les tables
	registerAllTables()
}

// registerAllTables enregistre toutes les tables du schéma
func registerAllTables() {
	// Table des utilisateurs
	registerTable("users", createUserTable())
	
	// Table des entreprises
	registerTable("companies", createCompanyTable())
	
	// Table des posts
	registerTable("posts", createPostTable())
	
	// Table des catégories
	registerTable("categories", createCategoryTable())
}

// registerTable ajoute une table au registre global
func registerTable(name string, builder *TableBuilder) {
	globalSchema.tables[name] = builder
	globalSchema.order = append(globalSchema.order, name)
	logging.Info.Printf("Table '%s' enregistrée dans le schéma", name)
}

// InitAllTables crée toutes les tables enregistrées dans la base de données
func InitAllTables(conn *Connection) error {
	logging.Info.Println("Initialisation de toutes les tables du schéma...")
	
	for _, tableName := range globalSchema.order {
		table := globalSchema.tables[tableName]
		
		logging.Info.Printf("Création de la table '%s'...", tableName)
		err := conn.CreateTable(table)
		if err != nil {
			return fmt.Errorf("erreur lors de la création de la table '%s': %v", tableName, err)
		}
		logging.Info.Printf("Table '%s' créée avec succès!", tableName)
	}
	
	logging.Info.Printf("Toutes les tables (%d) ont été créées avec succès!", len(globalSchema.order))
	return nil
}

// GetTable retourne une table spécifique du schéma
func GetTable(name string) (*TableBuilder, bool) {
	table, exists := globalSchema.tables[name]
	return table, exists
}

// GetAllTables retourne toutes les tables du schéma
func GetAllTables() map[string]*TableBuilder {
	return globalSchema.tables
}

// ListTables retourne la liste des noms de tables dans l'ordre d'enregistrement
func ListTables() []string {
	return globalSchema.order
}

// === DÉFINITIONS DES TABLES ===

// createUserTable crée la définition de la table users
func createUserTable() *TableBuilder {
	return NewTable("users").
		AddAttribute("name", String).NotNull().Build().
		AddAttribute("email", String).NotNull().Unique().Build().
		AddAttribute("password", String).NotNull().Build()
}

// createCompanyTable crée la définition de la table companies
func createCompanyTable() *TableBuilder {
	return NewTable("companies").
		AddAttribute("name", String).NotNull().Unique().Build().
		AddAttribute("description", String).Build().
		AddAttribute("employee_count", Integer).Build().
		AddAttribute("revenue", Float).Build().
		AddAttribute("is_public", Boolean).NotNull().Build()
}

// createPostTable crée la définition de la table posts
func createPostTable() *TableBuilder {
	return NewTable("posts").
		AddAttribute("title", String).NotNull().Build().
		AddAttribute("content", String).Build().
		AddAttribute("published", Boolean).Build()
}

// createCategoryTable crée la définition de la table categories
func createCategoryTable() *TableBuilder {
	return NewTable("categories").
		AddAttribute("slug", String).NotNull().Unique().Build().
		AddAttribute("display_name", String).NotNull().Build()
}
