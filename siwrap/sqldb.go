package siwrap

import (
	"context"
	"database/sql"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
)

// SqlDB is a wrapper of sql.DB
type SqlDB struct {
	db         *sql.DB
	sqlColumns []sicore.SqlColumn
}

// NewSqlDB returns SqlDB
func NewSqlDB(db *sql.DB, sc ...sicore.SqlColumn) *SqlDB {
	return &SqlDB{
		db:         db,
		sqlColumns: sc,
	}
}

// Begin begins a transaction
func (o *SqlDB) Begin() (*sql.Tx, error) {
	tx, err := o.db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (o *SqlDB) Prepare(query string) (*sql.Stmt, error) {
	return o.db.Prepare(query)
}

func (o *SqlDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return o.db.PrepareContext(ctx, query)
}

func (o *SqlDB) Query(query string, args ...any) (*sql.Rows, error) {
	return o.db.Query(query, args...)
}

func (o *SqlDB) QueryRow(query string, args ...any) *sql.Row {
	return o.db.QueryRow(query, args...)
}

func (o *SqlDB) Exec(query string, args ...any) (sql.Result, error) {
	return o.db.Exec(query, args...)
}

func (o *SqlDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return o.db.QueryContext(ctx, query, args...)
}
func (o *SqlDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return o.db.QueryRowContext(ctx, query, args...)
}

func (o *SqlDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return o.db.ExecContext(ctx, query, args...)
}

func (o *SqlDB) ExecRowsAffected(query string, args ...any) (int64, error) {
	res, err := o.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func (o *SqlDB) ExecContextRowsAffected(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := o.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// QueryMaps queries a database then scan resultset into output(slice of map)
func (o *SqlDB) QueryMaps(query string, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.db.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)

	return rs.Scan(rows, output, o.sqlColumns...)
}

// QueryStructs queries a database then scan resultset into output of any type
func (o *SqlDB) QueryStructs(query string, output any, args ...any) (int, error) {
	rows, err := o.db.Query(query, args...)
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

// QueryContenxtMaps queries a database with context then scan resultset into output(slice of map)
func (o *SqlDB) QueryContenxtMaps(ctx context.Context, query string, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)

	return rs.Scan(rows, output, o.sqlColumns...)
}

// QueryContextStructs queries a database with context then scan resultset into output of any type
func (o *SqlDB) QueryContextStructs(ctx context.Context, query string, output any, args ...any) (int, error) {
	rows, err := o.db.QueryContext(ctx, query, args...)
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
