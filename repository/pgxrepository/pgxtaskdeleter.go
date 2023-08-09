package pgxtaskrepo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxTaskDeleter struct {
	db *pgxpool.Pool
}

func NewPgxTaskDeleter(db *pgxpool.Pool) *PgxTaskDeleter {
	return &PgxTaskDeleter{
		db: db,
	}
}

func (d *PgxTaskDeleter) DeleteTask(ctx context.Context, taskId int64, userId int64) error {
	res, err := d.db.Exec(ctx, "UPDATE tasks SET deleted_at = $1 WHERE deleted_at IS NULL AND id = $2 AND user_id = $3", time.Now().Unix(), taskId, userId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("task not found")
	}
	return nil
}
