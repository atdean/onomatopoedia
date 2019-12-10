package models

type User struct {
	ID int `db:"user_id"`
	Username string `db:"username"`
	Email string `db:"email"`
	Password string `db:"password"`
}