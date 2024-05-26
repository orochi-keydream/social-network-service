package repository

import (
	"context"
	"database/sql"
)

type IExecutionContext interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type IConnectionFactory interface {
	GetMaster() *sql.DB
	GetSync() *sql.DB
	GetAsync() *sql.DB
}
