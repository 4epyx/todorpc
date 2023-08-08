package pgxtaskrepo

import (
	"context"
	"time"

	"github.com/4epyx/todorpc/pb"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxTaskCreator struct {
	db *pgxpool.Pool
}

func NewPgxTaskCreator(db *pgxpool.Pool) *PgxTaskCreator {
	return &PgxTaskCreator{
		db: db,
	}
}

func (c *PgxTaskCreator) CreateTask(ctx context.Context, task *pb.TaskRequest, userId int64) (*pb.Task, error) {
	res := &pb.Task{
		Title:       task.Title,
		Description: task.Description,
		Deadline:    task.Deadline,
	}
	if err := c.db.QueryRow(ctx, `INSERT INTO tasks (title, description, user_id, deadline, created_at) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`, task.Title, task.Description, userId, task.Deadline, time.Now().Unix()).
		Scan(&res.Id, &res.CreatedAt); err != nil {
		return nil, err
	}

	return res, nil
}
