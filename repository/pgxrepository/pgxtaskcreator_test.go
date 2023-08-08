package pgxtaskrepo_test

import (
	"context"
	"testing"
	"time"

	"github.com/4epyx/todorpc/pb"
	pgxrepo "github.com/4epyx/todorpc/repository/pgxrepository"
	"github.com/4epyx/todorpc/util/testutil"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestPgxCreator struct {
	suite.Suite
	db      *pgxpool.Pool
	creator *pgxrepo.PgxTaskCreator
	ctx     context.Context
}

func (t *TestPgxCreator) SetupTest() {
	var err error
	t.db, t.ctx, err = testutil.SetupTestDbConn()
	if err != nil {
		t.T().Fatal(err)
	}

	t.creator = pgxrepo.NewPgxTaskCreator(t.db)
}

func (t *TestPgxCreator) TestCreateValidTask() {
	validTask := &pb.TaskRequest{
		Title:       "Task 1",
		Description: "task1 description",
		Deadline:    time.Now().Add(time.Hour * 24).Unix(),
	}
	_, err := t.creator.CreateTask(t.ctx, validTask, 1)
	t.Nil(err)
}

func (t *TestPgxCreator) TestCreateInvalidTask() {
	invalidTask := &pb.TaskRequest{
		Description: "task2 description",
		Deadline:    time.Now().Add(time.Hour * 24).Unix(),
	}
	_, err := t.creator.CreateTask(t.ctx, invalidTask, 1)
	t.NotNil(err)
}

func (t *TestPgxCreator) TestCreateTaskWithSqlInjection() {
	injTask := &pb.TaskRequest{
		Title:       "Task 3",
		Description: "description'); DROP TABLE tasks;",
	}

	_, err := t.creator.CreateTask(t.ctx, injTask, 1)
	t.Nil(err)

	tables, err := testutil.GetAllTables(t.ctx, t.db)
	if err != nil {
		t.T().Fatal(err)
	}

	found := false
	for _, t := range tables {
		if t == "tasks" {
			found = true
			break
		}
	}

	t.True(found)
}

func (t *TestPgxCreator) TestCreateTaskWithSpecialChars() {
	scharsTask := &pb.TaskRequest{
		Title:       "Task 3",
		Description: "$8",
	}

	_, err := t.creator.CreateTask(t.ctx, scharsTask, 1)
	t.Nil(err)
}

func TestCreateTaskSuite(t *testing.T) {
	test := new(TestPgxCreator)
	suite.Run(t, test)
	test.db.Exec(context.Background(), "DROP TABLE IF EXISTS tasks")
}
