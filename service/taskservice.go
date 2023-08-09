package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/4epyx/todorpc/pb"
	"github.com/4epyx/todorpc/repository"
	"google.golang.org/grpc/metadata"
)

type TaskService struct {
	repo repository.TaskRepository
	pb.UnimplementedTaskServiceServer
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return TaskService{repo: repo}
}

func (t *TaskService) CreateTask(ctx context.Context, task *pb.TaskRequest) (*pb.Task, error) {
	id, err := t.getUserIdFromMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed when parsing user id: %w", err)
	}

	newTask, err := t.repo.CreateTask(ctx, task, id)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}

func (t *TaskService) GetAllTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.ShortTasks, error) {
	id, err := t.getUserIdFromMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed when parsing user id: %w", err)
	}

	tasks, err := t.repo.GetAllTasks(ctx, req, id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *TaskService) GetFullTaskInfo(ctx context.Context, taskId *pb.TaskId) (*pb.Task, error) {
	id, err := t.getUserIdFromMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed when parsing user id: %w", err)
	}

	task, err := t.repo.GetFullTaskInfo(ctx, taskId.Id, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}
func (t *TaskService) SetTaskCompleted(ctx context.Context, taskId *pb.TaskId) (*pb.Task, error) {
	id, err := t.getUserIdFromMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed when parsing user id: %w", err)
	}

	task, err := t.repo.SetTaskCompleted(ctx, taskId.Id, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *TaskService) SetTaskUncompleted(ctx context.Context, taskId *pb.TaskId) (*pb.Task, error) {
	id, err := t.getUserIdFromMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed when parsing user id: %w", err)
	}

	task, err := t.repo.SetTaskUncompleted(ctx, taskId.Id, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *TaskService) UpdateTask(ctx context.Context, task *pb.TaskToUpdate) (*pb.Task, error) {
	id, err := t.getUserIdFromMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed when parsing user id: %w", err)
	}

	updatedTask, err := t.repo.UpdateTask(ctx, task, id)
	if err != nil {
		return nil, err
	}
	return updatedTask, nil
}

func (t *TaskService) DeleteTask(ctx context.Context, taskId *pb.TaskId) (*pb.Empty, error) {
	id, err := t.getUserIdFromMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed when parsing user id: %w", err)
	}

	err = t.repo.DeleteTask(ctx, taskId.Id, id)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (t *TaskService) getUserIdFromMetadata(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("can not get metadata from context")
	}

	if len(md.Get("user_id")) == 0 {
		return 0, errors.New("can not get user id from metadata")
	}
	id, err := strconv.ParseInt(md.Get("user_id")[0], 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
