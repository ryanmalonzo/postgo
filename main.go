package main

import (
	"flag"
	"fmt"
	"postgo/db"
	"postgo/examples"
	"postgo/logging"

	_ "github.com/lib/pq"
)

func main() {
	var demo = flag.String("demo", "", "Demo to run: 'typed', or leave empty for basic")
	flag.Parse()

	switch *demo {
	case "typed":
		examples.DemoTypedInserts()
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

	// Initialisation automatique de toutes les tables du schéma
	err = db.InitAllTables(conn)
	if err != nil {
		panic(err)
	}
	logging.Info.Println("All schema tables initialized successfully!")
}
