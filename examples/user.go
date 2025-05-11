package examples

import (
	"postgo/db"
)

type User struct {
    db.Model
    Name     string `db:"name,not_null"`
    Email    string `db:"email,not_null,unique"`
    Password string `db:"password,not_null"`
}

func (u *User) TableName() string {
    return "users"
}