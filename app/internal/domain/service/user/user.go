package user

import (
	"context"
	"messenger-rest-api/app/internal/adapters/db/postgresql"
	"messenger-rest-api/app/internal/domain/entities"
	custom_error "messenger-rest-api/app/internal/errors"
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
		return "", custom_error.ErrUserDuplicate
	}

	user := entities.User{
		Username: dto.Username,
	}

	return u.storage.Create(ctx, user)
}
