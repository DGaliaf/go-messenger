package entities

import "time"

type User struct {
	ID        string    `json:"id,omitempty"` // User ID (UUID)
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
