package pgxtaskrepo

import (
	"context"
	"time"

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
	task := &pb.Task{}
	if err := u.db.QueryRow(ctx,
		"UPDATE tasks SET completed_at = $1 WHERE id = $2 AND user_id = $3 AND completed_at = 0 RETURNING title, description, deadline, created_at, completed_at",
		time.Now().Unix(), taskId, userId).
		Scan(&task.Title, &task.Description, &task.Deadline, &task.CreatedAt, &task.CompletedAt); err != nil {
		return nil, err
	}

	task.Id = taskId
	return task, nil
}

func (u *PgxTaskUpdater) SetTaskUncompleted(ctx context.Context, taskId int64, userId int64) (*pb.Task, error) {
	task := &pb.Task{}
	if err := u.db.QueryRow(ctx,
		"UPDATE tasks SET completed_at = $1 WHERE deleted_at IS NULL AND id = $2 AND user_id = $3 AND completed_at != 0 RETURNING title, description, deadline, created_at, completed_at",
		0, taskId, userId).
		Scan(&task.Title, &task.Description, &task.Deadline, &task.CreatedAt, &task.CompletedAt); err != nil {
		return nil, err
	}
	task.Id = taskId
	return task, nil
}

func (u *PgxTaskUpdater) UpdateTask(ctx context.Context, task *pb.TaskToUpdate, userId int64) (*pb.Task, error) {
	resultTask := &pb.Task{}
	if err := u.db.QueryRow(ctx,
		"UPDATE tasks SET title = $1, description = $2, deadline = $3 WHERE deleted_at IS NULL AND id = $4 AND user_id = $5 RETURNING title, description, deadline, created_at, completed_at",
		task.Title, task.Description, task.Deadline, task.Id, userId).
		Scan(&resultTask.Title, &resultTask.Description, &resultTask.Deadline, &resultTask.CreatedAt, &resultTask.CompletedAt); err != nil {
		return nil, err
	}
	resultTask.Id = task.Id
	return resultTask, nil
}
