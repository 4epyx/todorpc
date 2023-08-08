package pgxtaskrepo_test

import (
	"context"
	"sort"
	"sync"
	"testing"

	"github.com/4epyx/todorpc/pb"
	pgxrepo "github.com/4epyx/todorpc/repository/pgxrepository"
	"github.com/4epyx/todorpc/util/testutil"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestPgxReader struct {
	suite.Suite
	db     *pgxpool.Pool
	reader *pgxrepo.PgxTaskReader
	ctx    context.Context
	once   sync.Once
}

func (t *TestPgxReader) SetupTest() {
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
	t.reader = pgxrepo.NewPgxTaskReader(t.db)
}

func (t *TestPgxReader) TestGetAllTasksSortById() {
	req := &pb.GetTasksRequest{
		SortBy:        pb.SortBy_ID,
		ShowCompleted: true,
	}
	res, err := t.reader.GetAllTasks(t.ctx, req, 1)

	t.Nil(err)
	t.Equal(len(res.Tasks), 5)
	t.True(sort.SliceIsSorted(res.Tasks, func(i, j int) bool {
		return res.Tasks[i].Id < res.Tasks[j].Id
	}))
}

func (t *TestPgxReader) TestGetAllTasksSortByTitle() {
	req := &pb.GetTasksRequest{
		SortBy:        pb.SortBy_TITLE,
		ShowCompleted: true,
	}
	res, err := t.reader.GetAllTasks(t.ctx, req, 1)

	t.Nil(err)
	t.Equal(len(res.Tasks), 5)
	t.True(sort.SliceIsSorted(res.Tasks, func(i, j int) bool {
		return res.Tasks[i].Title < res.Tasks[j].Title
	}))
}

func (t *TestPgxReader) TestGetAllTasksSortByDeadline() {
	req := &pb.GetTasksRequest{
		SortBy:        pb.SortBy_DEADLINE,
		ShowCompleted: true,
	}
	res, err := t.reader.GetAllTasks(t.ctx, req, 1)

	t.Nil(err)
	t.Equal(len(res.Tasks), 5)
	t.True(sort.SliceIsSorted(res.Tasks, func(i, j int) bool {
		return res.Tasks[i].Deadline < res.Tasks[j].Deadline
	}))
}

func (t *TestPgxReader) TestGetAllTasksSortByCreatedAt() {
	req := &pb.GetTasksRequest{
		SortBy:        pb.SortBy_CREATED_AT,
		ShowCompleted: true,
	}
	res, err := t.reader.GetAllTasks(t.ctx, req, 1)

	t.Nil(err)
	t.Equal(len(res.Tasks), 5)
	t.True(sort.SliceIsSorted(res.Tasks, func(i, j int) bool {
		return res.Tasks[i].CreatedAt < res.Tasks[j].CreatedAt
	}))
}

func (t *TestPgxReader) TestGetUncompletedTasks() {
	req := &pb.GetTasksRequest{
		ShowCompleted: false,
	}
	res, err := t.reader.GetAllTasks(t.ctx, req, 1)

	t.Nil(err)
	t.Equal(len(res.Tasks), 1)
}

func TestPgxReaderSuit(t *testing.T) {
	test := new(TestPgxReader)
	suite.Run(t, test)
	test.db.Exec(context.Background(), "DROP TABLE IF EXISTS tasks")
}
