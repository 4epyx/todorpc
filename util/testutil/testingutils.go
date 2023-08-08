package testutil

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/4epyx/todorpc/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetAllTables(ctx context.Context, db *pgxpool.Pool) ([]string, error) {
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

func SetupTestDbConn() (*pgxpool.Pool, context.Context, error) {
	dbUrl, ok := os.LookupEnv("TEST_DB_URL")
	if !ok {
		return nil, nil, errors.New("test db url not found in environment variable")
	}

	ctx := context.Background()

	conn, err := db.ConnectToDB(ctx, dbUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("can not connect to db: %v", err)
	}

	if err := db.MigrateTaskTable(ctx, conn); err != nil {
		return nil, nil, fmt.Errorf("can not migrate to db: %v", err)
	}

	return conn, ctx, nil
}

func LoadDump(ctx context.Context, db *pgxpool.Pool) error {
	dump, err := readDumpFile()
	if err != nil {
		return err
	}

	if _, err := db.Exec(ctx, dump); err != nil {
		return err
	}

	return nil
}

func readDumpFile() (string, error) {
	data, err := os.ReadFile("../../util/testutil/dump.sql")

	if err != nil {
		return "", err
	}

	return string(data), nil
}
