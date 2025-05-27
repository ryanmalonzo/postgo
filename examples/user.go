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

func (u *User) TableName() string {
	return "users"
}
