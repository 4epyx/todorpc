package db_test

import (
	"context"
	"os"
	"testing"

	"github.com/4epyx/todorpc/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestMigrateToDb struct {
	suite.Suite
	conn *pgxpool.Pool
}

func getAllTables(ctx context.Context, db *pgxpool.Pool) ([]string, error) {
	res, err := db.Query(ctx, "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'")
	if err != nil {
		return nil, err
	}

	tables := make([]string, 0, 5)
	for res.Next() {
		var table string
		res.Scan(&table)
		tables = append(tables, table)
	}

	return tables, nil
}

func (t *TestMigrateToDb) SetupTest() {
	dbUrl, ok := os.LookupEnv("TEST_DB_URL")
	if !ok {
		t.T().Fatal("test db url not found in environment variable")
	}

	var err error
	t.conn, err = db.ConnectToDB(context.Background(), dbUrl)
	if err != nil {
		t.T().Fatal("failed to connect to db")
	}
}

func (t *TestMigrateToDb) TestMigrateTable() {
	err := db.MigrateTaskTable(context.Background(), t.conn)
	t.Nil(err)
	defer t.conn.Exec(context.Background(), "DROP TABLE tasks")

	tables, err := getAllTables(context.Background(), t.conn)
	t.Nil(err)

	found := false
	for _, t := range tables {
		if t == "tasks" {
			found = true
			break
		}
	}

	t.True(found)
}

func TestMigrateToDbSuite(t *testing.T) {
	suite.Run(t, new(TestMigrateToDb))
}
