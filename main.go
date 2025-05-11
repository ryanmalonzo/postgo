package main

import (
	"fmt"
	"postgo/db"
	"postgo/examples"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to the database
	conn, err := db.NewConnection("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// Create a table with a schema
	err = conn.CreateTable("users", examples.User{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database successfully!")
}
