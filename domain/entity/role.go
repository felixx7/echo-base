package entity

import "time"

// Role represents a role in the system
type Role struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RoleResponse represents role response
type RoleResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
