package dto

type CreateMessageDTO struct {
	ChatID   int    `json:"chat_id"`
	AuthorID string `json:"author_id"`
	Text     string `json:"text"`
}

type ShowMessageIdDTO struct {
	ID int `json:"message_id"`
}

type GetMessagesFromChatDTO struct {
	ChatID int `json:"chat_id"`
}
