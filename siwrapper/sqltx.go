package siwrap

import (
	"context"
	"database/sql"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
)

type SqlTx struct {
	tx         *sql.Tx
	sqlColumns []sicore.SqlColumn
}

func NewSqlTx(tx *sql.Tx, sc ...sicore.SqlColumn) *SqlTx {
	return &SqlTx{
		tx:         tx,
		sqlColumns: sc,
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

func (o *SqlTx) QueryRow(query string, args ...any) *sql.Row {
	return o.tx.QueryRow(query, args...)
}

func (o *SqlTx) Exec(query string, args ...any) (sql.Result, error) {
	return o.tx.Exec(query, args...)
}

func (o *SqlTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return o.tx.QueryContext(ctx, query, args...)
}

func (o *SqlTx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return o.tx.QueryRowContext(ctx, query, args...)
}

func (o *SqlTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return o.tx.ExecContext(ctx, query, args...)
}

func (o *SqlTx) QueryIntoMapSlice(query string, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.tx.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)

	return rs.Scan(rows, output, o.sqlColumns...)
}

func (o *SqlTx) QueryIntoAny(query string, output any, args ...any) (int, error) {
	rows, err := o.tx.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)

	list := make([]map[string]interface{}, 0)
	n, err := rs.Scan(rows, &list, o.sqlColumns...)
	if err != nil {
		return 0, err
	}

	// simple, not very ideal json unmarshal
	if err = siutils.DecodeAny(list[:n], output); err != nil {
		return 0, err
	}

	return n, nil
}

func (o *SqlTx) QueryContenxtIntoMapSlice(ctx context.Context, query string, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)

	return rs.Scan(rows, output, o.sqlColumns...)
}

func (o *SqlTx) QueryContextIntoAny(ctx context.Context, query string, output any, args ...any) (int, error) {
	rows, err := o.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)

	list := make([]map[string]interface{}, 0)
	n, err := rs.Scan(rows, &list, o.sqlColumns...)
	if err != nil {
		return 0, err
	}

	// simple, not very ideal json unmarshal
	if err = siutils.DecodeAny(list[:n], output); err != nil {
		return 0, err
	}

	return n, nil
}
