package repository

import (
	"context"

	"github.com/4epyx/todorpc/pb"
)

type TaskRepository interface {
	TaskReader
	TaskCreator
	TaskUpdater
	TaskDeleter
}

type TaskReader interface {
	GetAllTasks(context.Context, *pb.GetTasksRequest, int64 /*user id*/) (*pb.ShortTasks, error)
	GetFullTaskInfo(context.Context, int64 /*task id*/, int64 /*user id*/) (*pb.Task, error)
}

type TaskCreator interface {
	CreateTask(context.Context, *pb.TaskRequest, int64 /*user id*/) (*pb.Task, error)
}

type TaskUpdater interface {
	SetTaskCompleted(context.Context, int64 /*task id*/, int64 /*user id*/) (*pb.Task, error)
	SetTaskUncompleted(context.Context, int64 /*task id*/, int64 /*user id*/) (*pb.Task, error)
	UpdateTask(context.Context, *pb.TaskToUpdate, int64 /*user id*/) (*pb.Task, error)
}

type TaskDeleter interface {
	DeleteTask(context.Context, int64 /*task id*/, int64 /*user id*/) error
}
