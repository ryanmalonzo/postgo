package examples

import "postgo/db"

// CreateUserTable crée la définition de la table users en utilisant le builder pattern
func CreateUserTable() *db.TableBuilder {
	return db.NewTable("users").
		AddAttribute("name", db.String).NotNull().Build().
		AddAttribute("email", db.String).NotNull().Unique().Build().
		AddAttribute("password", db.String).NotNull().Build()
}
