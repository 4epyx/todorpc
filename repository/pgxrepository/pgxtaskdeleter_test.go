package pgxtaskrepo_test

import (
	"context"
	"sync"
	"testing"

	pgxrepo "github.com/4epyx/todorpc/repository/pgxrepository"
	"github.com/4epyx/todorpc/util/testutil"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestPgxDeleter struct {
	suite.Suite
	db      *pgxpool.Pool
	deleter *pgxrepo.PgxTaskDeleter
	ctx     context.Context
	once    sync.Once
}

func (t *TestPgxDeleter) SetupTest() {
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
	t.deleter = pgxrepo.NewPgxTaskDeleter(t.db)
}

func (t *TestPgxDeleter) TestDeleteTask() {
	err := t.deleter.DeleteTask(t.ctx, 8, 1)
	t.Nil(err)
}

func (t *TestPgxDeleter) TestDeleteBelongsToOtherTask() {
	err := t.deleter.DeleteTask(t.ctx, 7, 1)
	t.NotNil(err)
}

func (t *TestPgxDeleter) TestDeleteDeletedTask() {
	err := t.deleter.DeleteTask(t.ctx, 9, 1)
	t.NotNil(err)
}

func TestPgxDeleterSuit(t *testing.T) {
	test := new(TestPgxDeleter)
	suite.Run(t, test)
	test.db.Exec(context.Background(), "DROP TABLE IF EXISTS tasks")
}
