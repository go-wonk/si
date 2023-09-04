package sisql

import (
	"context"
	"database/sql"

	"github.com/go-wonk/si/v2/sicore"
)

type SqlStmt struct {
	stmt *sql.Stmt
	opts []sicore.RowScannerOption
}

func NewSqlStmt(stmt *sql.Stmt, opts ...sicore.RowScannerOption) *SqlStmt {
	return &SqlStmt{
		stmt: stmt,
		opts: opts,
	}
}

func (o *SqlStmt) QueryRow(args ...any) *sql.Row {
	return o.stmt.QueryRow(args...)
}

func (o *SqlStmt) QueryRowContext(ctx context.Context, args ...any) *sql.Row {
	return o.stmt.QueryRowContext(ctx, args...)
}

func (o *SqlStmt) Query(args ...any) (*sql.Rows, error) {
	return o.stmt.Query(args...)
}

func (o *SqlStmt) QueryContext(ctx context.Context, args ...any) (*sql.Rows, error) {
	return o.stmt.QueryContext(ctx, args...)
}

func (o *SqlStmt) Exec(args ...any) (sql.Result, error) {
	return o.stmt.Exec(args...)
}

func (o *SqlStmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	return o.stmt.ExecContext(ctx, args...)
}

func (o *SqlStmt) ExecRowsAffected(args ...any) (int64, error) {
	return o.ExecContextRowsAffected(context.Background(), args...)
}
func (o *SqlStmt) ExecContextRowsAffected(ctx context.Context, args ...any) (int64, error) {
	res, err := o.stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (o *SqlStmt) QueryMaps(output *[]map[string]interface{}, args ...any) (int, error) {
	return o.QueryContextMaps(context.Background(), output, args...)
}

func (o *SqlStmt) QueryContextMaps(ctx context.Context, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.stmt.QueryContext(ctx, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner(o.opts...)
	defer sicore.PutRowScanner(rs)

	return rs.ScanMapSlice(rows, output)
}

func (o *SqlStmt) QueryRowPrimary(output any, args ...any) error {
	return o.QueryRowContextPrimary(context.Background(), output, args...)
}

func (o *SqlStmt) QueryRowContextPrimary(ctx context.Context, output any, args ...any) error {
	row := o.stmt.QueryRowContext(ctx, args...)

	rs := sicore.GetRowScanner(o.opts...)
	defer sicore.PutRowScanner(rs)

	err := rs.ScanPrimary(row, output)
	if err != nil {
		return err
	}

	return nil
}

func (o *SqlStmt) QueryRowStruct(output any, args ...any) error {
	return o.QueryRowContextPrimary(context.Background(), output, args...)
}

func (o *SqlStmt) QueryRowContextStruct(ctx context.Context, output any, args ...any) error {
	rows, err := o.stmt.QueryContext(ctx, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner(o.opts...)
	defer sicore.PutRowScanner(rs)

	err = rs.ScanStruct(rows, output)
	if err != nil {
		return err
	}

	return nil
}

// QueryStructs queries a database then scan resultset into output of any type
func (o *SqlStmt) QueryStructs(output any, args ...any) (int, error) {
	return o.QueryContextStructs(context.Background(), output, args...)
}

// QueryContextStructs queries a database with context then scan resultset into output of any type
func (o *SqlStmt) QueryContextStructs(ctx context.Context, output any, args ...any) (int, error) {
	rows, err := o.stmt.QueryContext(ctx, args...)
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
