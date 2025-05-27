package examples

import "postgo/db"

// CreateCompanyTable démontre la création d'une table plus complexe
func CreateCompanyTable() *db.TableBuilder {
	return db.NewTable("companies").
		AddAttribute("name", db.String).NotNull().Unique().Build().
		AddAttribute("description", db.String).Build().
		AddAttribute("employee_count", db.Integer).Build().
		AddAttribute("revenue", db.Float).Build().
		AddAttribute("is_public", db.Boolean).NotNull().Build()
}

// CreatePostTable démontre une table simple
func CreatePostTable() *db.TableBuilder {
	return db.NewTable("posts").
		AddAttribute("title", db.String).NotNull().Build().
		AddAttribute("content", db.String).Build().
		AddAttribute("published", db.Boolean).Build()
}

// CreateCategoryTable démontre une table avec contraintes diverses
func CreateCategoryTable() *db.TableBuilder {
	return db.NewTable("categories").
		AddAttribute("slug", db.String).NotNull().Unique().Build().
		AddAttribute("display_name", db.String).NotNull().Build()
}
