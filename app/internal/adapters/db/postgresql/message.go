package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"messenger-rest-api/app/internal/domain/entities"
	"messenger-rest-api/app/internal/errors"
)

type MessageStorage struct {
	db *pgxpool.Pool
}

func NewMessageStorage(db *pgxpool.Pool) *MessageStorage {
	return &MessageStorage{db: db}
}

func (m MessageStorage) Create(ctx context.Context, message entities.Message) (int, error) {
	var id int

	acquire, err := m.db.Acquire(ctx)
	if err != nil {
		return 0, custom_error.ErrAcquireConnection
	}
	defer acquire.Release()

	sql := `INSERT INTO public.message(chat_id, author_id, text) VALUES ($1, $2, $3) RETURNING id`

	if err := acquire.QueryRow(ctx, sql, message.ChatID, message.AuthorID, message.Text).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, custom_error.ErrEntityNotFound
		}

		return 0, err
	}

	return id, nil
}
