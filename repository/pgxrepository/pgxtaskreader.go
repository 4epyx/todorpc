package pgxtaskrepo

import (
	"context"
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
	tasks := make([]*pb.ShortTask, 0)
	for i := 0; rows.Next(); i++ {
		task := &pb.ShortTask{}
		if err := rows.Scan(&task.Id, &task.Title, &task.Deadline, &task.CreatedAt, &task.CompletedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return &pb.ShortTasks{
		Tasks: tasks,
		Count: int32(len(tasks)),
	}, nil
}

func (r *PgxTaskReader) buildGetAllTasksQuery(showCompleted bool, sortBy pb.SortBy) string {
	query := &strings.Builder{}
	query.WriteString("SELECT id, title, deadline, created_at, completed_at FROM tasks WHERE deleted_at IS NULL AND user_id = $1 AND deleted_at IS NULL")
	if !showCompleted {
		query.WriteString(" AND completed_at = 0")
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
	task := &pb.Task{}
	if err := r.db.QueryRow(ctx, "SELECT id, title, description, deadline, created_at, completed_at FROM tasks WHERE deleted_at IS NULL AND id = $1 AND user_id = $2", taskId, userId).
		Scan(&task.Id, &task.Title, &task.Description, &task.Deadline, &task.CreatedAt, &task.CompletedAt); err != nil {
		return nil, err
	}

	return task, nil
}
