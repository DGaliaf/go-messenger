package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"messenger-rest-api/app/internal/domain/entities"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{db: db}
}

func (u UserStorage) Create(ctx context.Context, user entities.User) (string, error) {
	id := new(string)
	acquire, err := u.db.Acquire(ctx)
	if err != nil {
		// TODO: Error handling
		return "", err
	}
	defer acquire.Release()

	sql := `INSERT INTO public.user(username) VALUES ($1) RETURNING id`
	if err := acquire.QueryRow(ctx, sql, user.Username).Scan(id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// TODO: ErrNoRows
			return "", err
		}

		return "", err
	}

	return *id, nil
}

func (u UserStorage) IsExistsByUsername(ctx context.Context, username string) (bool, error) {
	count := new(int)
	sql := `SELECT COUNT(id) FROM public.user WHERE username=$1`

	if err := u.db.QueryRow(ctx, sql, username).Scan(count); err != nil {
		return false, err
	}

	if *count > 0 {
		return true, nil
	}

	return false, nil
}

func (u UserStorage) IsExistsByID(ctx context.Context, id string) (bool, error) {
	count := new(int)
	sql := `SELECT COUNT(id) FROM public.user WHERE id=$1`

	if err := u.db.QueryRow(ctx, sql, id).Scan(count); err != nil {
		return false, err
	}

	if *count > 0 {
		return true, nil
	}

	return false, nil
}
