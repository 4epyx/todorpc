package repository

import "github.com/4epyx/todorpc/pb"

type TaskRepository interface {
	TaskReader
	TaskCreator
	TaskUpdater
	TaskDeleter
}

type TaskReader interface {
	GetAllTasks(*pb.GetTasksRequest) *pb.ShortTasks
	GetFullTaskInfo(int64) *pb.Task
}

type TaskCreator interface {
	CreateTask(*pb.TaskRequest) pb.Task
}

type TaskUpdater interface {
	SetTaskCompleted(int64) *pb.Task
	SetTaskUncompleted(int64) *pb.Task
	UpdateTask(*pb.TaskToUpdate) *pb.Task
}

type TaskDeleter interface {
	DeleteTask(int64) *pb.Task
}
