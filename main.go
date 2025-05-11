package main

import (
	"postgo/db"
	"postgo/examples"
	"postgo/logging"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to the database
	conn, err := db.NewConnection("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	}
	logging.Info.Println("Connected to the database successfully!")

	defer conn.Close()

	// Create a table with a schema
	err = conn.CreateTable("users", examples.User{})
	if err != nil {
		panic(err)
	}
	logging.Info.Println("Table users created successfully!")

	// Test
	userMetadata, err := db.GetMetadata(examples.User{})
	if err != nil {
		panic(err)
	}
	logging.Info.Println("User metadata:", userMetadata)
}
