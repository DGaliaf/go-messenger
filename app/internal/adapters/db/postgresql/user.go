package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"messenger-rest-api/app/internal/domain/entities"
	"messenger-rest-api/app/internal/errors"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{db: db}
}

func (u UserStorage) Create(ctx context.Context, user entities.User) (string, error) {
	var id string
	acquire, err := u.db.Acquire(ctx)
	if err != nil {
		return "", custom_error.ErrEntityNotFound
	}
	defer acquire.Release()

	sql := `INSERT INTO public.user(username) VALUES ($1) RETURNING id`
	if err := acquire.QueryRow(ctx, sql, user.Username).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", custom_error.ErrEntityNotFound
		}

		return "", err
	}

	return id, nil
}

func (u UserStorage) IsExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int
	sql := `SELECT COUNT(id) FROM public.user WHERE username=$1`

	if err := u.db.QueryRow(ctx, sql, username).Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, custom_error.ErrEntityNotFound
		}

		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (u UserStorage) IsExistsByID(ctx context.Context, id string) (bool, error) {
	var count int

	sql := `SELECT COUNT(id) FROM public.user WHERE id=$1`
	if err := u.db.QueryRow(ctx, sql, id).Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, custom_error.ErrEntityNotFound
		}

		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (u UserStorage) IsExistsInChat(ctx context.Context, userID string, chatID int) (bool, error) {
	var count int

	sql := `SELECT COUNT(*) FROM public.user_chat WHERE user_id=$1 AND chat_id=$2`
	if err := u.db.QueryRow(ctx, sql, userID, chatID).Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, custom_error.ErrEntityNotFound
		}

		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
