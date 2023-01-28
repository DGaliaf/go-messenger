package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"messenger-rest-api/app/internal/domain/entities"
)

type ChatStorage struct {
	db *pgxpool.Pool
}

func NewChatStorage(db *pgxpool.Pool) *ChatStorage {
	return &ChatStorage{db: db}
}

func (c ChatStorage) Create(ctx context.Context, chat entities.Chat, usersID []string) (int, error) {
	id := new(int)

	acquire, err := c.db.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer acquire.Release()

	tx, err := acquire.Begin(ctx)
	if err != nil {
		// TODO: Error handling
		return 0, err
	}

	defer tx.Rollback(ctx)

	sql := `INSERT INTO public.chat(name) VALUES ($1) RETURNING id`
	if err := tx.QueryRow(ctx, sql, chat.Name).Scan(id); err != nil {
		//TODO: NoRows

		if errors.Is(err, pgx.ErrNoRows) {
			// TODO: Error handling
			return 0, err
		}

		return 0, err
	}

	for _, userID := range usersID {
		sql = `INSERT INTO public.user_chat(user_id, chat_id) VALUES ($1,$2)`
		res, err := tx.Exec(ctx, sql, userID, id)
		if err != nil {
			return 0, err
		}

		if res.RowsAffected() == 0 {
			// TODO: NoRowsAffected
			return 0, errors.New("no rows affected")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		if errors.Is(err, pgx.ErrTxClosed) {
			// TODO: Tx Already closed
			return 0, err
		}

		if errors.Is(err, pgx.ErrTxCommitRollback) {
			// TODO: Fail to commit changes
			return 0, err
		}

		return 0, err
	}

	return *id, nil
}

func (c ChatStorage) FindUserChats(ctx context.Context, id string) ([]entities.Chat, error) {
	acquire, err := c.db.Acquire(ctx)
	if err != nil {
		// TODO: Error handling
		return nil, err
	}

	sql := `SELECT chat_id FROM public.user_chat WHERE user_id=$1`
	chatQuery, err := acquire.Query(ctx, sql, id)
	if err != nil {
		// TODO: Error handling
		return nil, err
	}

	chatsIDs := make([]int, 0)
	for chatQuery.Next() {
		chatID := new(int)

		if err := chatQuery.Scan(chatID); err != nil {
			// TODO: Error handling
			return nil, err
		}

		chatsIDs = append(chatsIDs, *chatID)
	}

	chatQuery.Close()
	acquire.Release()

	chats := make([]entities.Chat, 0)
	for _, chatID := range chatsIDs {
		acquire, err = c.db.Acquire(ctx)
		if err != nil {
			// TODO: Error handling
			return nil, err
		}

		sql = `SELECT user_id FROM public.user_chat WHERE chat_id=$1`
		userQuery, err := acquire.Query(ctx, sql, chatID)
		if err != nil {
			// TODO: Error handling
			return nil, err
		}

		userIDs := make([]string, 0)
		for userQuery.Next() {
			userID := new(string)

			if err := userQuery.Scan(userID); err != nil {
				// TODO: Error handling
				return nil, err
			}

			userIDs = append(userIDs, *userID)
		}

		userQuery.Close()
		acquire.Release()

		acquire, err = c.db.Acquire(ctx)

		users := make([]entities.User, 0)
		for _, userID := range userIDs {
			user := entities.User{}
			sql = `SELECT * FROM public.user WHERE id=$1`

			if err := acquire.QueryRow(ctx, sql, userID).Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
				// TODO: Error handling
				return nil, err
			}

			users = append(users, user)
		}

		chat := entities.Chat{}
		sql = `SELECT * FROM public.chat WHERE id=$1`
		if err := acquire.QueryRow(ctx, sql, chatID).Scan(&chat.ID, &chat.Name, &chat.CreatedAt); err != nil {
			// TODO: Error handling
			return nil, err
		}

		acquire.Release()

		chat.Users = users
		chats = append(chats, chat)
	}

	return chats, nil
}

func (c ChatStorage) IsExists(ctx context.Context, name string) (bool, error) {
	count := new(int)
	sql := `SELECT COUNT(id) FROM public.chat WHERE name=$1`

	if err := c.db.QueryRow(ctx, sql, name).Scan(count); err != nil {
		return false, err
	}

	if *count > 0 {
		return true, nil
	}

	return false, nil
}
