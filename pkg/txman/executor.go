package txman

import (
	"context"
	"database/sql"
)

type Executor interface {
	Prepare(query string) (stmt *sql.Stmt, err error)
	PrepareContext(ctx context.Context, query string) (stmt *sql.Stmt, err error)

	Exec(query string, args ...any) (result sql.Result, err error)
	ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error)

	Query(query string, args ...any) (rows *sql.Rows, err error)
	QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)

	QueryRow(query string, args ...any) (row *sql.Row)
	QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row)
}

var _ Executor = (*sql.DB)(nil)

var _ Executor = (*sql.Tx)(nil)
