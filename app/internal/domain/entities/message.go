package entities

import "time"

type Message struct {
	Id        int       `json:"message_id,omitempty"`
	ChatID    int       `json:"chat_id,omitempty"`
	AuthorID  string    `json:"author_id,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
