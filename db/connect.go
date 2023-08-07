package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectToDB(ctx context.Context, connectionURL string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(ctx, connectionURL)
}
