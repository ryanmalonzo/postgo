package main

import (
	"fmt"
	"postgo/examples"
	"postgo/functions"

	_ "github.com/lib/pq"
)

func main() {
	// Create the database first
	err := functions.CreateDatabase("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	} 

	// Then connect to the database
	db, err := functions.Connect("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	}
	
	// Create a table with a schema
	err = functions.CreateTable(db, "users", examples.User{})
	if err != nil {
		panic(err)
	}

	defer db.Close()
	fmt.Println("Connected to the database successfully!")
}
