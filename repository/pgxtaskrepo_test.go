package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/4epyx/todorpc/db"
	"github.com/4epyx/todorpc/pb"
	"github.com/4epyx/todorpc/repository"
	"github.com/4epyx/todorpc/util/testutil"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestPgxRepository struct {
	suite.Suite
	db   *pgxpool.Pool
	repo *repository.PgxTaskRepository
	ctx  context.Context
}

func (t *TestPgxRepository) SetupTest() {
	dbUrl, ok := os.LookupEnv("TEST_DB_URL")
	if !ok {
		t.T().Fatal("test db url not found in environment variable")
	}

	t.ctx = context.Background()

	var err error
	t.db, err = db.ConnectToDB(t.ctx, dbUrl)
	if err != nil {
		t.T().Fatalf("can not connect to db: %v", err)
	}

	if err := db.MigrateTaskTable(t.ctx, t.db); err != nil {
		t.T().Fatalf("can not migrate to db: %v", err)
	}
	t.repo = repository.NewTaskRepository(t.db)
}

func (t *TestPgxRepository) TestCreateValidTask() {
	validTask := &pb.TaskRequest{
		Title:       "Task 1",
		Description: "task1 description",
		Deadline:    time.Now().Add(time.Hour * 24).Unix(),
	}
	_, err := t.repo.CreateTask(t.ctx, validTask, 1)
	t.Nil(err)
}

func (t *TestPgxRepository) TestCreateInvalidTask() {
	invalidTask := &pb.TaskRequest{
		Description: "task2 description",
		Deadline:    time.Now().Add(time.Hour * 24).Unix(),
	}
	_, err := t.repo.CreateTask(t.ctx, invalidTask, 1)
	t.NotNil(err)
}

func (t *TestPgxRepository) TestCreateTaskWithSqlInjection() {
	injTask := &pb.TaskRequest{
		Title:       "Task 3",
		Description: "description'); DROP TABLE tasks;",
	}

	_, err := t.repo.CreateTask(t.ctx, injTask, 1)
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

func TestCreateTaskSuite(t *testing.T) {
	test := new(TestPgxRepository)
	suite.Run(t, test)
	test.db.Exec(context.Background(), "DROP TABLE IF EXISTS tasks")
}
