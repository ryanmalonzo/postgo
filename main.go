package main

import (
	"flag"
	"fmt"
	"log"
	"postgo/db"
	"postgo/db/query"
	"postgo/examples"
	"postgo/logging"

	_ "github.com/lib/pq"
)

func main() {
	var demo = flag.String("demo", "", "Demo to run: 'builder', 'full', or leave empty for basic")
	flag.Parse()

	switch *demo {
	case "builder":
		runExampleBuilder()
	case "full":
		runDemo()
	default:
		runBasicDemo()
	}
}

func runBasicDemo() {
	fmt.Println("=== Démonstration de l'ORM PostgreSQL ===")

	// Connexion à la base de données
	conn, err := db.NewConnection("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	}
	logging.Info.Println("Connected to the database successfully!")

	defer conn.Close()

	// Création d'une table avec le builder pattern
	// La table est définie explicitement avec ses attributs et contraintes
	userTable := examples.CreateUserTable()
	err = conn.CreateTable(userTable)
	if err != nil {
		panic(err)
	}
	logging.Info.Println("Table users created successfully!")
}

func runDemo() {
	fmt.Println("=== Démonstration du Builder Pattern pour l'ORM PostgreSQL ===")
	fmt.Println()

	// 1. Génération du SQL sans connexion à la base
	fmt.Println("1. Génération du SQL pour différentes tables:")
	fmt.Println()

	tables := []struct {
		name    string
		builder *db.TableBuilder
	}{
		{"Users", examples.CreateUserTable()},
		{"Companies", examples.CreateCompanyTable()},
		{"Posts", examples.CreatePostTable()},
		{"Categories", examples.CreateCategoryTable()},
	}

	for _, table := range tables {
		fmt.Printf("Table %s:\n", table.name)
		fmt.Printf("%s\n\n", table.builder.BuildSQL())
	}

	// 2. Connexion à la base et création réelle des tables
	fmt.Println("2. Connexion à la base de données et création des tables:")
	fmt.Println()

	conn, err := db.NewConnection("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		log.Printf("Erreur de connexion (ignorée pour la démo): %v", err)
		fmt.Println("Note: La base de données n'est pas disponible, mais le SQL généré est valide.")
		fmt.Println()
		return
	}
	defer conn.Close()

	logging.Info.Println("Connecté à la base de données avec succès!")

	// Création de toutes les tables
	for _, table := range tables {
		err = conn.CreateTable(table.builder)
		if err != nil {
			log.Printf("Erreur lors de la création de la table %s: %v", table.name, err)
		} else {
			logging.Info.Printf("Table %s créée avec succès!", table.name)
		}
	}

	// 3. Test des requêtes INSERT et SELECT
	fmt.Println("3. Test des requêtes INSERT et SELECT:")
	fmt.Println()

	insertQuery := query.NewInsertQuery("users").
		AddColumn("id").AddValue(1).
		AddColumn("name").AddValue("John").
		AddColumn("email").AddValue("john@example.com").
		AddColumn("age").AddValue(25).
		AddColumn("isactive").AddValue(true)

	_, err = insertQuery.Execute(conn.GetDatabase())
	if err != nil {
		log.Printf("Erreur lors de l'exécution de l'INSERT: %v", err)
	} else {
		fmt.Println("INSERT réussi!")
	}

	selectQuery := query.NewSelectQuery("users").
		AddColumn("id").
		AddColumn("name").
		Where("id = 1")

	rows, err := selectQuery.Execute(conn.GetDatabase())
	if err != nil {
		log.Printf("Erreur lors de l'exécution du SELECT: %v", err)
		return
	}
	defer rows.Close()

	fmt.Println("Résultats du SELECT:")
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Printf("Erreur lors du scan: %v", err)
			continue
		}
		fmt.Printf("User trouvé: ID=%d, Name=%s\n", id, name)
	}

	updateQuery := query.NewUpdateQuery("users").
		AddColumn("name").AddValue("John Does").
		Where("id = 1")

	err = updateQuery.Execute(conn.GetDatabase())
	if err != nil {
		log.Printf("Erreur lors de l'exécution de l'UPDATE: %v", err)
	}

	deleteQuery := query.NewDeleteQuery("users").
		AddCondition("id = 1")

	err = deleteQuery.Execute(conn.GetDatabase())
	if err != nil {
		log.Printf("Erreur lors de l'exécution du DELETE: %v", err)
	}

	fmt.Println()
	fmt.Println("=== Fin de la démonstration ===")
}

func runExampleBuilder() {
	fmt.Println("=== Démonstration du Builder Pattern ===")
	fmt.Println()

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

	fmt.Println()
	fmt.Println("=== Fin de la démonstration Builder ===")
}
