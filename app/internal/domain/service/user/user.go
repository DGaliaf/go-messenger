package user

import (
	"context"
	"errors"
	"messenger-rest-api/app/internal/adapters/db/postgresql"
	"messenger-rest-api/app/internal/domain/entities"
)

type Service struct {
	storage *postgresql.UserStorage
}

func NewUserService(storage *postgresql.UserStorage) *Service {
	return &Service{storage: storage}
}

func (u Service) Create(ctx context.Context, dto CreateUserDTO) (string, error) {
	exists, err := u.storage.IsExistsByUsername(ctx, dto.Username)
	if err != nil {
		return "", err
	}

	if exists {
		// TODO: ERROR Handling
		return "", errors.New("duplicate")
	}

	user := entities.User{
		Username: dto.Username,
	}

	return u.storage.Create(ctx, user)
}
