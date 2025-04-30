package main

import (
	"fmt"
	"postgo/functions"

	_ "github.com/lib/pq"
)

func main() {
	// Create the database first
	err := functions.CreateDatabase("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	}
	fmt.Println("Database created or already exists!")

	// Then connect to the database
	db, err := functions.Connect("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	fmt.Println("Connected to the database successfully!")
}
