package dto

type CreateChatDTO struct {
	Name    string   `json:"name"`  // Chat Name
	UsersID []string `json:"users"` // User IDs (UUIDs)
}

type ShowChatIdDTO struct {
	ID int `json:"chat_id"`
}
