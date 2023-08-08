package pgxtaskrepo

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxTaskRepository struct {
	*PgxTaskCreator
	*PgxTaskReader
	*PgxTaskUpdater
	*PgxTaskDeleter
}

func NewTaskRepository(db *pgxpool.Pool) *PgxTaskRepository {
	return &PgxTaskRepository{
		PgxTaskCreator: NewPgxTaskCreator(db),
		PgxTaskReader:  NewPgxTaskReader(db),
		PgxTaskUpdater: NewPgxTaskUpdater(db),
		PgxTaskDeleter: NewPgxTaskDeleter(db),
	}
}
