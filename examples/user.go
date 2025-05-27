package examples

import (
	"postgo/db"
)

type User struct {
	db.BaseModel
	Name     string `db:"not_null"`
	Email    string `db:"not_null,unique"`
	Password string `db:"not_null"`
}

// TableName retourne le nom de la table pour ce modèle.
// Cette méthode implémente l'interface db.Model.
func (u *User) TableName() string {
	return "users"
}
