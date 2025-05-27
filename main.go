package main

import (
	"flag"
	"postgo/db"
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
