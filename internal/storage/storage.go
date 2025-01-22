package storage

import (
	"context"
	"database/sql"
	"log/slog"
)

type storage struct {
	db DBTX
	l  *slog.Logger
}

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func New(db DBTX, cacheLogger *slog.Logger) *storage {
	return &storage{
		db: db,
		l:  cacheLogger,
	}
}
