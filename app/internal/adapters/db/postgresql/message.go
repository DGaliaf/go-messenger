package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"messenger-rest-api/app/internal/domain/entities"
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
		return 0, err
	}
	defer acquire.Release()

	sql := `INSERT INTO public.message(chat_id, author_id, text) VALUES ($1, $2, $3) RETURNING id`

	if err := acquire.QueryRow(ctx, sql, message.ChatID, message.AuthorID, message.Text).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// TODO: Error handling
			return 0, err
		}

		return 0, err
	}

	return id, nil
}
