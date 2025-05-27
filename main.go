package main

import (
	"postgo/db"
	"postgo/examples"
	"postgo/logging"

	_ "github.com/lib/pq"
)

func main() {
	// Connexion à la base de données
	conn, err := db.NewConnection("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	}
	logging.Info.Println("Connected to the database successfully!")

	defer conn.Close()

	// Création d'une table avec un schéma basé sur le modèle User
	// Cette opération utilise la réflexion pour analyser la structure du modèle
	// et générer automatiquement le SQL de création de table
	err = conn.CreateTable(&examples.User{})
	if err != nil {
		panic(err)
	}
	logging.Info.Println("Table users created successfully!")
}
