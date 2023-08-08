package pgxtaskrepo

import (
	"context"
	"errors"
	"strings"

	"github.com/4epyx/todorpc/pb"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxTaskReader struct {
	db *pgxpool.Pool
}

func NewPgxTaskReader(db *pgxpool.Pool) *PgxTaskReader {
	return &PgxTaskReader{
		db: db,
	}
}

func (r *PgxTaskReader) GetAllTasks(ctx context.Context, req *pb.GetTasksRequest, userId int64) (*pb.ShortTasks, error) {
	query := r.buildGetAllTasksQuery(req.ShowCompleted, req.SortBy)
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	tasks := make([]*pb.ShortTask, rows.CommandTag().RowsAffected())
	for i := 0; rows.Next(); i++ {
		if err := rows.Scan(&tasks[i].Id, &tasks[i].Title, &tasks[i].Deadline, &tasks[i].CreatedAt, &tasks[i].CompletedAt); err != nil {
			return nil, err
		}
	}

	return &pb.ShortTasks{
		Tasks: tasks,
		Count: int32(len(tasks)),
	}, nil
}

func (r *PgxTaskReader) buildGetAllTasksQuery(showCompleted bool, sortBy pb.SortBy) string {
	query := &strings.Builder{}
	query.WriteString("SELECT id, title, deadline, created_at, completed_at FROM tasks WHERE user_id = $1")
	if !showCompleted {
		query.WriteString(" AND completed_at IS NULL")
	}

	query.WriteString(" ORDER BY ")
	switch sortBy {
	case pb.SortBy_ID:
		query.WriteString(" id")
	case pb.SortBy_TITLE:
		query.WriteString(" title")
	case pb.SortBy_CREATED_AT:
		query.WriteString(" created_at")
	case pb.SortBy_DEADLINE:
		query.WriteString(" deadline")
	default:
		query.WriteString(" created_at")
	}

	return query.String()
}

func (r *PgxTaskReader) GetFullTaskInfo(ctx context.Context, taskId int64, userId int64) (*pb.Task, error) {
	return nil, errors.New("not implemented")
}
