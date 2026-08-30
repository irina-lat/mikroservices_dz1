package model

import "time"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	SessionUUID string `json:"session_uuid"`
}

type WhoamiResponse struct {
	Session *SessionInfo `json:"session"`
	User    *User        `json:"user"`
}

type SessionInfo struct {
	UUID      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ExpiresAt string `json:"expires_at"`
}

// Session — модель для бизнес-логики
type Session struct {
	UUID      string    `json:"uuid"`
	UserUUID  string    `json:"user_uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}