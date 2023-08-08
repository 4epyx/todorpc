package repository

import (
	"context"
	"time"

	"github.com/4epyx/todorpc/pb"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxTaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *PgxTaskRepository {
	return &PgxTaskRepository{
		db: db,
	}
}

func (r *PgxTaskRepository) CreateTask(ctx context.Context, task *pb.TaskRequest, userId int64) (*pb.Task, error) {
	res := &pb.Task{
		Title:       task.Title,
		Description: task.Description,
		Deadline:    task.Deadline,
	}
	if err := r.db.QueryRow(ctx, `INSERT INTO tasks (title, description, user_id, deadline, created_at) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`, task.Title, task.Description, userId, task.Deadline, time.Now().Unix()).
		Scan(&res.Id, &res.CreatedAt); err != nil {
		return nil, err
	}

	return res, nil
}
