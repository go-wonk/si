package siwrap

import (
	"context"
	"database/sql"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
)

type SqlStmt struct {
	stmt       *sql.Stmt
	sqlColumns []sicore.SqlColumn
}

func NewSqlStmt(stmt *sql.Stmt, sc ...sicore.SqlColumn) *SqlStmt {
	return &SqlStmt{
		stmt:       stmt,
		sqlColumns: sc,
	}
}

func (o *SqlStmt) Query(args ...any) (*sql.Rows, error) {
	return o.stmt.Query(args...)
}

func (o *SqlStmt) Exec(args ...any) (sql.Result, error) {
	return o.stmt.Exec(args...)
}

func (o *SqlStmt) QueryRow(args ...any) *sql.Row {
	return o.stmt.QueryRow(args...)
}

func (o *SqlStmt) QueryContext(ctx context.Context, args ...any) (*sql.Rows, error) {
	return o.stmt.QueryContext(ctx, args...)
}

func (o *SqlStmt) QueryRowContext(ctx context.Context, args ...any) *sql.Row {
	return o.stmt.QueryRowContext(ctx, args...)
}

func (o *SqlStmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	return o.stmt.ExecContext(ctx, args...)
}

func (o *SqlStmt) ExecRowsAffected(args ...any) (int64, error) {
	res, err := o.stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func (o *SqlStmt) ExecContextRowsAffected(ctx context.Context, args ...any) (int64, error) {
	res, err := o.stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (o *SqlStmt) QueryMaps(output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.stmt.Query(args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)

	return rs.Scan(rows, output, o.sqlColumns...)
}

func (o *SqlStmt) QueryStructs(output any, args ...any) (int, error) {
	rows, err := o.stmt.Query(args...)
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

func (o *SqlStmt) QueryContextMaps(ctx context.Context, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.stmt.QueryContext(ctx, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)

	return rs.Scan(rows, output, o.sqlColumns...)
}

func (o *SqlStmt) QueryContextStructs(ctx context.Context, output any, args ...any) (int, error) {
	rows, err := o.stmt.QueryContext(ctx, args...)
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
