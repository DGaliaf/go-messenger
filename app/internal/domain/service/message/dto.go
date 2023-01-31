package message

type CreateMessageDTO struct {
	ChatID   int    `json:"chat_id"`
	AuthorID string `json:"author_id"`
	Text     string `json:"text"`
}

type GetMessagesFromChatDTO struct {
	ChatID int `json:"chat_id"`
}
