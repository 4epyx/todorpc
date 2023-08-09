package service_test

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/4epyx/todorpc/pb"
	pgxtaskrepo "github.com/4epyx/todorpc/repository/pgxrepository"
	"github.com/4epyx/todorpc/service"
	"github.com/4epyx/todorpc/util/testutil"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
)

type TestTaskService struct {
	suite.Suite
	ctx     context.Context
	db      *pgxpool.Pool
	service service.TaskService
	once    sync.Once
}

func (t *TestTaskService) SetupTest() {
	var err error
	t.db, _, err = testutil.SetupTestDbConn()
	if err != nil {
		t.T().Fatal(err)
	}
	t.ctx = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{}))
	t.once.Do(func() {
		if err := testutil.LoadDump(context.Background(), t.db); err != nil {
			t.T().Fatal(err)
		}
	})
	repo := pgxtaskrepo.NewTaskRepository(t.db)
	t.service = service.NewTaskService(repo)
}

func (t *TestTaskService) TestCreateValidTask() {
	validTask := &pb.TaskRequest{
		Title:       "Task 1",
		Description: "task1 description",
		Deadline:    time.Now().Add(time.Hour * 24).Unix(),
	}

	task, err := t.service.CreateTask(interceptorMock(t.ctx, 10), validTask)
	if !t.Nil(err) {
		t.T().Log(err)
	}

	t.Equal(task.Title, validTask.Title)
	t.Equal(task.Description, validTask.Description)
	t.Equal(task.Deadline, validTask.Deadline)
	t.NotEqual(task.CreatedAt, 0)
	t.NotEqual(task.Id, 0)
}

func (t *TestTaskService) TestCreateInvalidTask() {
	invalidTask := &pb.TaskRequest{
		Description: "task2 description",
		Deadline:    time.Now().Add(time.Hour * 24).Unix(),
	}
	_, err := t.service.CreateTask(interceptorMock(t.ctx, 10), invalidTask)
	t.NotNil(err)
}

func (t *TestTaskService) TestGetAllTasks() {
	req := &pb.GetTasksRequest{
		SortBy:        pb.SortBy_ID,
		ShowCompleted: true,
	}
	res, err := t.service.GetAllTasks(interceptorMock(t.ctx, 1), req)

	t.Nil(err)
	t.Equal(len(res.Tasks), 6)
	t.True(sort.SliceIsSorted(res.Tasks, func(i, j int) bool {
		return res.Tasks[i].Id < res.Tasks[j].Id
	}))
}

func (t *TestTaskService) TestGetUncompletedTasks() {
	req := &pb.GetTasksRequest{
		ShowCompleted: false,
	}
	res, err := t.service.GetAllTasks(interceptorMock(t.ctx, 1), req)

	t.Nil(err)
	t.Equal(len(res.Tasks), 2)
}

func (t *TestTaskService) TestGetNoTasks() {
	req := &pb.GetTasksRequest{
		ShowCompleted: true,
	}
	res, err := t.service.GetAllTasks(interceptorMock(t.ctx, -10), req)

	t.Nil(err)
	t.Equal(len(res.Tasks), 0)
}

func (t *TestTaskService) TestGetOneTask() {
	expected := &pb.Task{
		Id:          1,
		Title:       "full task",
		Description: "description",
		Deadline:    1691998920,
		CreatedAt:   1691491320,
		CompletedAt: 1691492320,
	}

	got, err := t.service.GetFullTaskInfo(interceptorMock(t.ctx, 1), &pb.TaskId{Id: 1})
	t.Nil(err)
	t.Equal(expected, got)
}

func (t *TestTaskService) TestGetNotExistingTask() {
	_, err := t.service.GetFullTaskInfo(interceptorMock(t.ctx, 1), &pb.TaskId{Id: 10})
	t.NotNil(err)
}

func (t *TestTaskService) TestSetTaskCompleted() {
	expected := &pb.Task{
		Id:          5,
		Title:       "uncompleted task",
		Description: "description",
		Deadline:    1691998920,
		CreatedAt:   1691491319,
	}

	res, err := t.service.SetTaskCompleted(interceptorMock(t.ctx, 1), &pb.TaskId{Id: 5})
	t.Nil(err)

	t.Equal(res.Id, expected.Id)
	t.Equal(res.Title, expected.Title)
	t.Equal(res.Description, expected.Description)
	t.Equal(res.Deadline, expected.Deadline)
	t.Equal(res.CreatedAt, expected.CreatedAt)
	t.NotEqual(res.CompletedAt, 0)
}

func (t *TestTaskService) TestSetCompletedBelongingToOtherTask() {
	_, err := t.service.SetTaskCompleted(interceptorMock(t.ctx, 1), &pb.TaskId{Id: 6})
	t.NotNil(err)
}

func (t *TestTaskService) TestSetTaskUncompleted() {
	expected := &pb.Task{
		Id:          1,
		Title:       "full task",
		Description: "description",
		Deadline:    1691998920,
		CreatedAt:   1691491320,
		CompletedAt: 0,
	}

	res, err := t.service.SetTaskUncompleted(interceptorMock(t.ctx, 1), &pb.TaskId{Id: 1})
	t.Nil(err)

	t.Equal(expected, res)
}

func (t *TestTaskService) TestSetUncompletedBelongingToOtherTask() {
	_, err := t.service.SetTaskUncompleted(interceptorMock(t.ctx, 1), &pb.TaskId{Id: 7})
	t.NotNil(err)
}

func (t *TestTaskService) TestUpdateTask() {
	upd := &pb.TaskToUpdate{
		Id:          3,
		Title:       "new task3 title",
		Description: "new task3 description",
		Deadline:    1691998920,
	}
	expected := &pb.Task{
		Id:          3,
		Title:       "new task3 title",
		Description: "new task3 description",
		Deadline:    1691998920,
		CreatedAt:   1691491320,
		CompletedAt: 1691492320,
	}

	res, err := t.service.UpdateTask(interceptorMock(t.ctx, 1), upd)
	t.Nil(err)
	t.Equal(expected, res)
}

func (t *TestTaskService) TestUpdateBelongingToOtherTask() {
	otherUserTask := &pb.TaskToUpdate{
		Id:          6,
		Title:       "new task title",
		Description: "new task description",
		Deadline:    1691998920,
	}

	_, err := t.service.UpdateTask(interceptorMock(t.ctx, 1), otherUserTask)
	t.NotNil(err)
}

func interceptorMock(ctx context.Context, userId int64) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"user_id": strconv.FormatInt(userId, 10)}))
}

func (t *TestTaskService) TestDeleteTask() {
	_, err := t.service.DeleteTask(interceptorMock(t.ctx, 1), &pb.TaskId{Id: 8})
	t.Nil(err)
}

func (t *TestTaskService) TestDeleteBelongsToOtherTask() {
	_, err := t.service.DeleteTask(interceptorMock(t.ctx, 1), &pb.TaskId{Id: 7})
	t.NotNil(err)
}

func TestTaskServiceSuite(t *testing.T) {
	test := new(TestTaskService)
	suite.Run(t, test)
	test.db.Exec(context.Background(), "DROP TABLE IF EXISTS tasks")
}
