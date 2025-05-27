package main

import (
	"fmt"
	"postgo/db"
	"postgo/examples"
)

func runExampleBuilder() {
	// Exemple d'utilisation du nouveau builder pattern pour créer des tables

	// Création d'une table simple avec des attributs
	usersTable := db.NewTable("users").
		AddAttribute("name", db.String).NotNull().Build().
		AddAttribute("email", db.String).NotNull().Unique().Build().
		AddAttribute("password", db.String).NotNull().Build()

	fmt.Println("SQL généré pour la table users:")
	fmt.Println(usersTable.BuildSQL())
	fmt.Println()

	// Création d'une autre table d'exemple
	productsTable := db.NewTable("products").
		AddAttribute("name", db.String).NotNull().Build().
		AddAttribute("price", db.Float).NotNull().Build().
		AddAttribute("in_stock", db.Boolean).Build()

	fmt.Println("SQL généré pour la table products:")
	fmt.Println(productsTable.BuildSQL())
	fmt.Println()

	// Test des nouvelles tables avancées
	companyTable := examples.CreateCompanyTable()
	fmt.Println("SQL généré pour la table companies:")
	fmt.Println(companyTable.BuildSQL())
	fmt.Println()

	postTable := examples.CreatePostTable()
	fmt.Println("SQL généré pour la table posts:")
	fmt.Println(postTable.BuildSQL())
	fmt.Println()

	categoryTable := examples.CreateCategoryTable()
	fmt.Println("SQL généré pour la table categories:")
	fmt.Println(categoryTable.BuildSQL())
	fmt.Println()

	// Création d'une table avec seulement l'ID (démonstration de l'ID auto-incrémenté obligatoire)
	simpleTable := db.NewTable("simple")

	fmt.Println("SQL généré pour la table simple (seulement ID):")
	fmt.Println(simpleTable.BuildSQL())
}
