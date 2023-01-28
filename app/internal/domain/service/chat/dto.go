package chat

type CreateChatDTO struct {
	Name    string   `json:"name"`
	UsersID []string `json:"users"`
}
