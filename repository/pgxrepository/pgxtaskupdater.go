package pgxtaskrepo

import (
	"context"
	"errors"

	"github.com/4epyx/todorpc/pb"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxTaskUpdater struct {
	db *pgxpool.Pool
}

func NewPgxTaskUpdater(db *pgxpool.Pool) *PgxTaskUpdater {
	return &PgxTaskUpdater{
		db: db,
	}
}

func (u *PgxTaskUpdater) SetTaskCompleted(ctx context.Context, taskId int64, userId int64) (*pb.Task, error) {
	return nil, errors.New("not implemented")
}
func (u *PgxTaskUpdater) SetTaskUncompleted(ctx context.Context, taskId int64, userId int64) (*pb.Task, error) {
	return nil, errors.New("not implemented")
}
func (u *PgxTaskUpdater) UpdateTask(ctx context.Context, task *pb.TaskToUpdate, userId int64) (*pb.Task, error) {
	return nil, errors.New("not implemented")
}
