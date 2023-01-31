package message

import (
	"context"
	"errors"
	"messenger-rest-api/app/internal/adapters/db/postgresql"
	"messenger-rest-api/app/internal/domain/entities"
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
	// TOOD: Mapping
	exists, err := s.userStorage.IsExistsByID(ctx, dto.AuthorID)
	if err != nil {
		// TODO: Error handling
		return 0, err
	}

	if !exists {
		// TODO: Error handling
		return 0, errors.New("user does not exist")
	}

	exists, err = s.chatStorage.IsExistsByID(ctx, dto.ChatID)
	if err != nil {
		// TODO: Error handling
		return 0, err
	}

	if !exists {
		// TODO: Error handling
		return 0, errors.New("user does not exist")
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
		// TODO: Error handling
		return nil, err
	}

	if !exists {
		// TODO: Error handling
		return nil, errors.New("chat does not exist")
	}

	return s.chatStorage.GetMessages(ctx, dto.ChatID)
}
