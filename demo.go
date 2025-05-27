package main

import (
	"fmt"
	"log"
	"postgo/db"
	"postgo/examples"
	"postgo/logging"

	_ "github.com/lib/pq"
)

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

	fmt.Println()
	fmt.Println("=== Fin de la démonstration ===")
}
