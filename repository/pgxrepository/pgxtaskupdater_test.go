package pgxtaskrepo_test

import (
	"context"
	"sync"
	"testing"

	"github.com/4epyx/todorpc/pb"
	pgxrepo "github.com/4epyx/todorpc/repository/pgxrepository"
	"github.com/4epyx/todorpc/util/testutil"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestPgxUpdater struct {
	suite.Suite
	db      *pgxpool.Pool
	updater *pgxrepo.PgxTaskUpdater
	ctx     context.Context
	once    sync.Once
}

func (t *TestPgxUpdater) SetupTest() {
	var err error
	t.db, t.ctx, err = testutil.SetupTestDbConn()
	if err != nil {
		t.T().Fatal(err)
	}

	t.once.Do(func() {
		if err := testutil.LoadDump(t.ctx, t.db); err != nil {
			t.T().Fatal(err)
		}
	})
	t.updater = pgxrepo.NewPgxTaskUpdater(t.db)
}

func (t *TestPgxUpdater) TestSetTaskCompleted() {
	expected := &pb.Task{
		Id:          5,
		Title:       "uncompleted task",
		Description: "description",
		Deadline:    1691998920,
		CreatedAt:   1691491319,
	}

	res, err := t.updater.SetTaskCompleted(t.ctx, 5, 1)
	t.Nil(err)

	t.Equal(res.Id, expected.Id)
	t.Equal(res.Title, expected.Title)
	t.Equal(res.Description, expected.Description)
	t.Equal(res.Deadline, expected.Deadline)
	t.Equal(res.CreatedAt, expected.CreatedAt)
	t.NotEqual(res.CompletedAt, 0)
}

func (t *TestPgxUpdater) TestSetCompletedTaskCompleted() {
	_, err := t.updater.SetTaskCompleted(t.ctx, 1, 1)
	t.NotNil(err)
}

func (t *TestPgxUpdater) TestSetCompletedBelongingToOtherTask() {
	_, err := t.updater.SetTaskCompleted(t.ctx, 6, 1)
	t.NotNil(err)
}

func (t *TestPgxUpdater) TestSetTaskUncompleted() {
	expected := &pb.Task{
		Id:          1,
		Title:       "full task",
		Description: "description",
		Deadline:    1691998920,
		CreatedAt:   1691491320,
		CompletedAt: 0,
	}

	res, err := t.updater.SetTaskUncompleted(t.ctx, 1, 1)
	t.Nil(err)

	t.Equal(expected, res)
}

func (t *TestPgxUpdater) TestSetUncompletedTaskUncompleted() {
	_, err := t.updater.SetTaskUncompleted(t.ctx, 8, 1)
	t.NotNil(err)
}

func (t *TestPgxUpdater) TestSetUncompletedBelongingToOtherTask() {
	_, err := t.updater.SetTaskUncompleted(t.ctx, 7, 1)
	t.NotNil(err)
}

func (t *TestPgxUpdater) TestUpdateTask() {
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

	res, err := t.updater.UpdateTask(t.ctx, upd, 1)
	t.Nil(err)
	t.Equal(expected, res)
}

func (t *TestPgxUpdater) TestUpdateNotExistingTask() {
	notExistingTask := &pb.TaskToUpdate{
		Id:          -1,
		Title:       "new task title",
		Description: "new task description",
		Deadline:    1691998920,
	}

	_, err := t.updater.UpdateTask(t.ctx, notExistingTask, 1)
	t.NotNil(err)
}

func (t *TestPgxUpdater) TestUpdateBelongingToOtherTask() {
	notExistingTask := &pb.TaskToUpdate{
		Id:          6,
		Title:       "new task title",
		Description: "new task description",
		Deadline:    1691998920,
	}

	_, err := t.updater.UpdateTask(t.ctx, notExistingTask, 1)
	t.NotNil(err)
}

func TestTaskUpdaterSuite(t *testing.T) {
	test := new(TestPgxUpdater)
	suite.Run(t, test)
	test.db.Exec(context.Background(), "DROP TABLE IF EXISTS tasks")
}
