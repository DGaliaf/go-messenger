package dto

type CreateChatDTO struct {
	Name    string   `json:"name"`
	UsersID []string `json:"users"`
}

type ShowChatIdDTO struct {
	ID int `json:"chat_id"`
}
