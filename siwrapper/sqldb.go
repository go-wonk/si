package siwrapper

import (
	"context"
	"database/sql"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
)

type SqlDB struct {
	db         *sql.DB
	sqlColumns []sicore.SqlColumn
}

func NewSqlDB(db *sql.DB) *SqlDB {
	return &SqlDB{
		db: db,
	}
}

func (o *SqlDB) AddSqlColumn(sc ...sicore.SqlColumn) {
	o.sqlColumns = sc
}

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

func (o *SqlDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return o.db.ExecContext(ctx, query, args...)
}

func (o *SqlDB) QueryIntoMap(query string, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.db.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)
	for _, c := range o.sqlColumns {
		c.SetType(rs)
	}

	scannedRow, columns, err := rs.ScanTypes(rows)
	if err != nil {
		return 0, err
	}

	n := 0
	for rows.Next() {
		err = rows.Scan(scannedRow...)
		if err != nil {
			return 0, err
		}

		m := make(map[string]interface{})
		rs.ScanValuesIntoMap(columns, scannedRow, &m)
		*output = append(*output, m)
		n++
	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (o *SqlDB) QueryIntoAny(query string, output any, args ...any) (int, error) {
	rows, err := o.db.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner()
	defer sicore.PutRowScanner(rs)
	for _, c := range o.sqlColumns {
		c.SetType(rs)
	}

	list := make([]map[string]interface{}, 0)
	scannedRow, columns, err := rs.ScanTypes(rows)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		err = rows.Scan(scannedRow...)
		if err != nil {
			return 0, err
		}

		m := make(map[string]interface{})
		rs.ScanValuesIntoMap(columns, scannedRow, &m)
		list = append(list, m)

	}

	err = rows.Err()
	if err != nil {
		return 0, err
	}

	if len(list) == 0 {
		return 0, nil
	}

	// simple, not very ideal json unmarshal
	if err = siutils.DecodeAny(list, output); err != nil {
		return 0, err
	}

	return 1, nil
}
