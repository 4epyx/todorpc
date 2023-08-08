package pgxtaskrepo

import (
	"context"
	"errors"

	"github.com/4epyx/todorpc/pb"
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

func (d *PgxTaskDeleter) DeleteTask(ctx context.Context, taskId int64, userId int64) (*pb.Task, error) {
	return nil, errors.New("not implemented")
}
