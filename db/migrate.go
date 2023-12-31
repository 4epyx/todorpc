package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

const createTaskTableQuery = `CREATE TABLE IF NOT EXISTS tasks (
	id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
	title TEXT NOT NULL CHECK (char_length(title) > 3),
	description TEXT,
	user_id BIGINT NOT NULL,
	deadline BIGINT NOT NULL DEFAULT 0,
	created_at BIGINT NOT NULL,
	completed_at BIGINT NOT NULL DEFAULT 0,
	deleted_at BIGINT DEFAULT NULL
)`

func MigrateTaskTable(ctx context.Context, conn *pgxpool.Pool) error {
	_, err := conn.Exec(ctx, createTaskTableQuery)
	return err
}
