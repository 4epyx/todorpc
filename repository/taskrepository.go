package repository

import "github.com/4epyx/todorpc/pb"

type TaskRepository interface {
	TaskReader
	TaskCreator
	TaskUpdater
	TaskDeleter
}

type TaskReader interface {
	GetAllTasks(*pb.GetTasksRequest, int64 /*user id*/) *pb.ShortTasks
	GetFullTaskInfo(int64 /*task id*/, int64 /*user id*/) *pb.Task
}

type TaskCreator interface {
	CreateTask(*pb.TaskRequest, int64 /*user id*/) pb.Task
}

type TaskUpdater interface {
	SetTaskCompleted(int64 /*task id*/, int64 /*user id*/) *pb.Task
	SetTaskUncompleted(int64 /*task id*/, int64 /*user id*/) *pb.Task
	UpdateTask(*pb.TaskToUpdate, int64 /*user id*/) *pb.Task
}

type TaskDeleter interface {
	DeleteTask(int64 /*task id*/, int64 /*user id*/) *pb.Task
}
