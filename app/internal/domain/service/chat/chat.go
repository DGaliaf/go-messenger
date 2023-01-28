package chat

import (
	"context"
	"errors"
	"messenger-rest-api/app/internal/adapters/db/postgresql"
	"messenger-rest-api/app/internal/domain/entities"
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
	exists, err := c.chatStorage.IsExists(ctx, dto.Name)
	if err != nil {
		return 0, err
	}

	if exists {
		// TODO: ERROR Handling
		return 0, errors.New("duplicate")
	}

	if len(dto.UsersID) < 2 {
		// TODO: ERROR Handling
		return 0, errors.New("not enough users")
	}

	if containDuplicates(dto.UsersID) {
		// TODO: ERROR Handling
		return 0, errors.New("contain users with the same ID")
	}

	chat := entities.Chat{
		Name: dto.Name,
	}

	return c.chatStorage.Create(ctx, chat, dto.UsersID)
}

func (c Service) FindUserChats(ctx context.Context, id string) ([]entities.Chat, error) {
	exists, err := c.userStorage.IsExistsByID(ctx, id)
	if err != nil {
		// TODO: Error Handling
		return nil, err
	}

	if !exists {
		// TODO: Error Handling
		return nil, errors.New("user does not exists")
	}

	return c.chatStorage.FindUserChats(ctx, id)
}

func containDuplicates[T comparable](data []T) bool {
	counts := make(map[T]int)

	for _, d := range data {
		counts[d]++
	}

	for _, val := range counts {
		if val > 1 {
			return true
		}
	}

	return false
}
