package model

import "time"

type User struct {
	UUID      string    `json:"uuid"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // не возвращаем в JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}