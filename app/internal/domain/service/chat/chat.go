package chat

import (
	"context"
	"messenger-rest-api/app/internal/adapters/db/postgresql"
	"messenger-rest-api/app/internal/domain/entities"
	"messenger-rest-api/app/internal/errors"
	"messenger-rest-api/app/pkg/utils"
)

type Service struct {
	chatStorage *postgresql.ChatStorage
	userStorage *postgresql.UserStorage
}

func NewChatService(chatStorage *postgresql.ChatStorage, userStorage *postgresql.UserStorage) *Service {
	return &Service{
		chatStorage: chatStorage,
		userStorage: userStorage,
	}
}

func (c Service) Create(ctx context.Context, dto CreateChatDTO) (int, error) {
	exists, err := c.chatStorage.IsExistsByName(ctx, dto.Name)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, custom_error.ErrChatDuplicate
	}

	if len(dto.UsersID) < 2 {
		return 0, custom_error.ErrNotEnoughUsers
	}

	if utils.ContainDuplicates(dto.UsersID) {
		return 0, custom_error.ErrUserAlreadyInChat
	}

	chat := entities.Chat{
		Name: dto.Name,
	}

	return c.chatStorage.Create(ctx, chat, dto.UsersID)
}

func (c Service) FindUserChats(ctx context.Context, id string) ([]entities.Chat, error) {
	exists, err := c.userStorage.IsExistsByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, custom_error.ErrUserNotExist
	}

	return c.chatStorage.FindUserChats(ctx, id)
}
