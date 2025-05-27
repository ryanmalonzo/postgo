package examples

import (
	"postgo/db"
)

// User représente un utilisateur dans le système.
// Cette structure utilise l'embedding de BaseModel pour hériter
// automatiquement du champ ID avec ses contraintes (primary_key, auto_increment).
// Les tags `db` définissent les contraintes SQL pour chaque champ.
type User struct {
	db.BaseModel                    // Héritage du modèle de base (ID, etc.)
	Name     string `db:"not_null"` // Nom obligatoire
	Email    string `db:"not_null,unique"` // Email obligatoire et unique
	Password string `db:"not_null"` // Mot de passe obligatoire
}

// TableName retourne le nom de la table pour ce modèle.
// Cette méthode implémente l'interface db.Model.
func (u *User) TableName() string {
	return "users"
}
