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

type ChatStorage struct {
	db *pgxpool.Pool
}

func NewChatStorage(db *pgxpool.Pool) *ChatStorage {
	return &ChatStorage{db: db}
}

func (c ChatStorage) Create(ctx context.Context, chat entities.Chat, usersID []string) (int, error) {
	var id int

	acquire, err := c.db.Acquire(ctx)
	if err != nil {
		return 0, custom_error.ErrAcquireConnection
	}
	defer acquire.Release()

	tx, err := acquire.Begin(ctx)
	if err != nil {
		return 0, custom_error.ErrCreateTransaction
	}

	defer tx.Rollback(ctx)

	sql := `INSERT INTO public.chat(name) VALUES ($1) RETURNING id`
	if err := tx.QueryRow(ctx, sql, chat.Name).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, custom_error.ErrEntityNotFound
		}

		return 0, err
	}

	for _, userID := range usersID {
		sql = `INSERT INTO public.user_chat(user_id, chat_id) VALUES ($1,$2)`
		res, err := tx.Exec(ctx, sql, userID, id)
		if err != nil {
			return 0, custom_error.ErrSQLExecution
		}

		if res.RowsAffected() == 0 {
			return 0, custom_error.ErrNoRowsAffected
		}
	}

	if err := tx.Commit(ctx); err != nil {
		if errors.Is(err, pgx.ErrTxClosed) {
			return 0, custom_error.ErrTransactionClosed
		}

		if errors.Is(err, pgx.ErrTxCommitRollback) {
			return 0, custom_error.ErrCommitTransaction
		}

		return 0, err
	}

	return id, nil
}

func (c ChatStorage) FindUserChats(ctx context.Context, id string) ([]entities.Chat, error) {
	acquire, err := c.db.Acquire(ctx)
	if err != nil {
		return nil, custom_error.ErrAcquireConnection
	}

	sql := `SELECT chat_id FROM public.user_chat WHERE user_id=$1`
	chatQuery, err := acquire.Query(ctx, sql, id)
	if err != nil {
		return nil, custom_error.ErrSQLExecution
	}

	chatsIDs := make([]int, 0)
	for chatQuery.Next() {
		var chatID int

		if err := chatQuery.Scan(&chatID); err != nil {
			return nil, err
		}

		chatsIDs = append(chatsIDs, chatID)
	}

	chatQuery.Close()
	acquire.Release()

	chats := make([]entities.Chat, 0)
	for _, chatID := range chatsIDs {
		acquire, err = c.db.Acquire(ctx)
		if err != nil {
			return nil, custom_error.ErrAcquireConnection
		}

		sql = `SELECT user_id FROM public.user_chat WHERE chat_id=$1`
		userQuery, err := acquire.Query(ctx, sql, chatID)
		if err != nil {
			return nil, custom_error.ErrSQLExecution
		}

		userIDs := make([]string, 0)
		for userQuery.Next() {
			var userID string

			if err := userQuery.Scan(&userID); err != nil {
				return nil, err
			}

			userIDs = append(userIDs, userID)
		}

		userQuery.Close()
		acquire.Release()

		acquire, err = c.db.Acquire(ctx)
		if err != nil {
			return nil, custom_error.ErrAcquireConnection
		}

		users := make([]entities.User, 0)
		for _, userID := range userIDs {
			user := entities.User{}
			sql = `SELECT * FROM public.user WHERE id=$1`

			if err := acquire.QueryRow(ctx, sql, userID).Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, custom_error.ErrEntityNotFound
				}

				return nil, err
			}

			users = append(users, user)
		}

		chat := entities.Chat{}
		sql = `SELECT * FROM public.chat WHERE id=$1`
		if err := acquire.QueryRow(ctx, sql, chatID).Scan(&chat.ID, &chat.Name, &chat.CreatedAt); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, custom_error.ErrEntityNotFound
			}

			return nil, err
		}

		acquire.Release()

		chat.Users = users
		chats = append(chats, chat)
	}

	return chats, nil
}

func (c ChatStorage) GetMessages(ctx context.Context, id int) ([]entities.Message, error) {
	acquire, err := c.db.Acquire(ctx)
	if err != nil {
		return nil, custom_error.ErrAcquireConnection
	}
	defer acquire.Release()

	sql := `SELECT * FROM public.message WHERE chat_id=$1 ORDER BY created_at DESC`
	query, err := acquire.Query(ctx, sql, id)
	if err != nil {
		return nil, custom_error.ErrSQLExecution
	}

	messages := make([]entities.Message, 0)
	for query.Next() {
		message := entities.Message{}

		if err := query.Scan(&message.Id, &message.ChatID, &message.AuthorID, &message.Text, &message.CreatedAt); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (c ChatStorage) IsExistsByID(ctx context.Context, id int) (bool, error) {
	var count int
	sql := `SELECT COUNT(id) FROM public.chat WHERE id=$1`

	if err := c.db.QueryRow(ctx, sql, id).Scan(&count); err != nil {
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

func (c ChatStorage) IsExistsByName(ctx context.Context, name string) (bool, error) {
	var count int

	sql := `SELECT COUNT(id) FROM public.chat WHERE name=$1`
	if err := c.db.QueryRow(ctx, sql, name).Scan(&count); err != nil {
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
