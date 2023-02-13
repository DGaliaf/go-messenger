package message

import (
	"context"
	"messenger-rest-api/app/internal/adapters/db/postgresql"
	"messenger-rest-api/app/internal/domain/entities"
	custom_error "messenger-rest-api/app/internal/errors"
)

type Service struct {
	messageStorage *postgresql.MessageStorage
	userStorage    *postgresql.UserStorage
	chatStorage    *postgresql.ChatStorage
}

func NewMessageService(messageStorage *postgresql.MessageStorage, userStorage *postgresql.UserStorage, chatStorage *postgresql.ChatStorage) *Service {
	return &Service{
		messageStorage: messageStorage,
		userStorage:    userStorage,
		chatStorage:    chatStorage,
	}
}

func (s Service) SendMessage(ctx context.Context, dto CreateMessageDTO) (int, error) {
	exists, err := s.userStorage.IsExistsByID(ctx, dto.AuthorID)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, custom_error.ErrUserNotExist
	}

	exists, err = s.userStorage.IsExistsInChat(ctx, dto.AuthorID, dto.ChatID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, custom_error.ErrUserNotInChat
	}

	exists, err = s.chatStorage.IsExistsByID(ctx, dto.ChatID)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, custom_error.ErrChatNotExist
	}

	message := entities.Message{
		ChatID:   dto.ChatID,
		AuthorID: dto.AuthorID,
		Text:     dto.Text,
	}

	return s.messageStorage.Create(ctx, message)
}

func (s Service) GetMessagesFromChatByID(ctx context.Context, dto GetMessagesFromChatDTO) ([]entities.Message, error) {
	exists, err := s.chatStorage.IsExistsByID(ctx, dto.ChatID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, custom_error.ErrChatNotExist
	}

	return s.chatStorage.GetMessages(ctx, dto.ChatID)
}
