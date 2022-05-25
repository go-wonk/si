package siwrapper

import (
	"context"
	"database/sql"
)

type SqlTx struct {
	tx *sql.Tx
}

func NewSqlTx(tx *sql.Tx) *SqlTx {
	return &SqlTx{
		tx: tx,
	}
}

func (o *SqlTx) Commit() error {
	return o.tx.Commit()
}

func (o *SqlTx) Rollback() error {
	return o.tx.Rollback()
}

func (o *SqlTx) Prepare(query string) (*sql.Stmt, error) {
	return o.tx.Prepare(query)
}

func (o *SqlTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return o.tx.PrepareContext(ctx, query)
}

func (o *SqlTx) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return o.tx.Query(query, args...)
}

func (o *SqlTx) Exec(query string, args ...any) (sql.Result, error) {
	return o.tx.Exec(query, args...)
}

func (o *SqlTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return o.tx.QueryContext(ctx, query, args...)
}

func (o *SqlTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return o.tx.ExecContext(ctx, query, args...)
}
