package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type pgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func NewPgConfig(username string, password string, host string, port string, database string) *pgConfig {
	return &pgConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

func NewClient(ctx context.Context, maxAttempts int, maxDelay time.Duration, cfg *pgConfig) (*pgxpool.Pool, error) {
	var client *pgxpool.Pool

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port, cfg.Database,
	)

	err := DoWithAttempts(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pgxConfig, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return err
		}
		pgxConfig.MaxConns = 60000

		client, err = pgxpool.NewWithConfig(ctx, pgxConfig)
		if err != nil {
			return err
		}

		if err := client.Ping(ctx); err != nil {
			return err
		}

		initSQL := `BEGIN;

CREATE SCHEMA IF NOT EXISTS public;
SET search_path TO public;

CREATE TABLE IF NOT EXISTS public.user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS public.chat (
    id SERIAL PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS public.user_chat (
    user_id UUID REFERENCES public.user(id),
    chat_id INT REFERENCES public.chat(id)
);

CREATE TABLE IF NOT EXISTS public.message (
    id SERIAL PRIMARY KEY NOT NULL,
    chat_id INT REFERENCES public.chat(id),
    author_id UUID REFERENCES public.user(id),
    text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

END;`
		if _, err = client.Exec(ctx, initSQL); err != nil {
			return err
		}

		return nil
	}, maxAttempts, maxDelay)

	return client, err
}

func DoWithAttempts(fn func() error, maxAttempts int, delay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--

			continue
		}

		return nil
	}

	return err
}
