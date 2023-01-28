package dto

type CreateUserDTO struct {
	Username string `json:"username"`
}

type GetUserDTO struct {
	UserID string `json:"user_id"`
}

type ShowUserIdDTO struct {
	UserID string `json:"user_id"`
}
