package entities

import "time"

type Chat struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Users     []User    `json:"users,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
