package testutil

import (
	"context"

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
