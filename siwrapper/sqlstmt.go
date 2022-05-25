package siwrapper

import (
	"context"
	"database/sql"
)

type SqlStmt struct {
	stmt *sql.Stmt
}

func NewSqlStmt(stmt *sql.Stmt) *SqlStmt {
	return &SqlStmt{
		stmt: stmt,
	}
}

func (o *SqlStmt) Query(args ...any) (*sql.Rows, error) {
	return o.stmt.Query(args...)
}

func (o *SqlStmt) Exec(args ...any) (sql.Result, error) {
	return o.stmt.Exec(args...)
}

func (o *SqlStmt) QueryContext(ctx context.Context, args ...any) (*sql.Rows, error) {
	return o.stmt.QueryContext(ctx, args...)
}

func (o *SqlStmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	return o.stmt.ExecContext(ctx, args...)
}
