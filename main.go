package main

import (
	"fmt"
	"postgo/functions"

	_ "github.com/lib/pq"
)

func main() {
	db, err := functions.Connect("localhost", 5432, "postgo", "postgo", "postgo")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	fmt.Println("Connected to the database successfully!")
}
