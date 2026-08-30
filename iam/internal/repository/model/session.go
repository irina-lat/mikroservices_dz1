package model

import "time"

type Session struct {
	UUID      string    `json:"uuid"`
	UserUUID  string    `json:"user_uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}