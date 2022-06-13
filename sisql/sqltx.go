package sisql

import (
	"context"
	"database/sql"

	"github.com/go-wonk/si/sicore"
)

type SqlTx struct {
	tx   *sql.Tx
	opts []sicore.RowScannerOption
}

func newSqlTx(tx *sql.Tx, opts ...sicore.RowScannerOption) *SqlTx {
	return &SqlTx{
		tx:   tx,
		opts: opts,
	}
}

func (o *SqlTx) Reset(tx *sql.Tx) {
	o.tx = tx
	o.opts = o.opts[:0]
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

func (o *SqlTx) ExecRowsAffected(query string, args ...any) (int64, error) {
	res, err := o.tx.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func (o *SqlTx) ExecContextRowsAffected(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := o.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (o *SqlTx) QueryMaps(query string, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.tx.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner(o.opts...)
	defer sicore.PutRowScanner(rs)

	return rs.ScanMapSlice(rows, output)
}

func (o *SqlTx) QueryStructs(query string, output any, args ...any) (int, error) {
	rows, err := o.tx.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner(o.opts...)
	defer sicore.PutRowScanner(rs)

	n, err := rs.ScanStructs(rows, output)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (o *SqlTx) QueryContextMaps(ctx context.Context, query string, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner(o.opts...)
	defer sicore.PutRowScanner(rs)

	return rs.ScanMapSlice(rows, output)
}

func (o *SqlTx) QueryContextStructs(ctx context.Context, query string, output any, args ...any) (int, error) {
	rows, err := o.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner(o.opts...)
	defer sicore.PutRowScanner(rs)

	n, err := rs.ScanStructs(rows, output)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (o *SqlTx) WithTagKey(key string) *SqlTx {
	o.opts = append(o.opts, sicore.WithTagKey(key))
	return o
}
