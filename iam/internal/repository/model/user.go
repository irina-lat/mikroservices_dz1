package model

import "time"

type User struct {
	UUID      string    `db:"uuid"`
	Login     string    `db:"login"`
	Email     string    `db:"email"`
	Password  string    `db:"password"` // хешированный пароль
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}