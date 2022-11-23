package sisql

import (
	"context"
	"database/sql"
)

type Querier interface {
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryMaps(query string, output *[]map[string]interface{}, args ...any) (int, error)
	QueryContextMaps(ctx context.Context, query string, output *[]map[string]interface{}, args ...any) (int, error)
	QueryRowPrimary(query string, output any, args ...any) error
	QueryRowContextPrimary(ctx context.Context, query string, output any, args ...any) error
	QueryRowStruct(query string, output any, args ...any) error
	QueryRowContextStruct(ctx context.Context, query string, output any, args ...any) error
	QueryStructs(query string, output any, args ...any) (int, error)
	QueryContextStructs(ctx context.Context, query string, output any, args ...any) (int, error)
}

type Executor interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	ExecRowsAffected(query string, args ...any) (int64, error)
	ExecContextRowsAffected(ctx context.Context, query string, args ...any) (int64, error)
}
