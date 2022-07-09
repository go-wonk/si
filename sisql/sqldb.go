package sisql

import (
	"context"
	"database/sql"

	"github.com/go-wonk/si/sicore"
)

func Open(driverName string, dataSourceName string) (*SqlDB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return NewSqlDB(db), nil
}

// SqlDB is a wrapper of sql.DB
type SqlDB struct {
	db   *sql.DB
	opts []sicore.RowScannerOption
}

// NewSqlDB returns SqlDB
func NewSqlDB(db *sql.DB, opts ...SqlOption) *SqlDB {
	sqldb := &SqlDB{
		db: db,
		// opts: opts,
	}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o.apply(sqldb)
	}

	return sqldb
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

func (o *SqlDB) QueryRow(query string, args ...any) *sql.Row {
	return o.db.QueryRow(query, args...)
}

func (o *SqlDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return o.db.QueryRowContext(ctx, query, args...)
}

func (o *SqlDB) Query(query string, args ...any) (*sql.Rows, error) {
	return o.db.Query(query, args...)
}

func (o *SqlDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return o.db.QueryContext(ctx, query, args...)
}

func (o *SqlDB) Exec(query string, args ...any) (sql.Result, error) {
	return o.db.Exec(query, args...)
}

func (o *SqlDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return o.db.ExecContext(ctx, query, args...)
}

// ExecRowsAffected executes query and returns number of affected rows.
func (o *SqlDB) ExecRowsAffected(query string, args ...any) (int64, error) {
	return o.ExecContextRowsAffected(context.Background(), query, args...)
}

// ExecContextRowsAffected executes query and returns number of affected rows.
func (o *SqlDB) ExecContextRowsAffected(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := o.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// QueryMaps queries a database then scan resultset into output(slice of map)
func (o *SqlDB) QueryMaps(query string, output *[]map[string]interface{}, args ...any) (int, error) {
	return o.QueryContextMaps(context.Background(), query, output, args...)
}

// QueryContextMaps queries a database with context then scan resultset into output(slice of map)
func (o *SqlDB) QueryContextMaps(ctx context.Context, query string, output *[]map[string]interface{}, args ...any) (int, error) {
	rows, err := o.db.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rs := sicore.GetRowScanner(o.opts...)
	defer sicore.PutRowScanner(rs)

	return rs.ScanMapSlice(rows, output)
}

func (o *SqlDB) QueryRowPrimary(query string, output any, args ...any) error {
	return o.QueryRowContextPrimary(context.Background(), query, output, args...)
}

func (o *SqlDB) QueryRowContextPrimary(ctx context.Context, query string, output any, args ...any) error {
	row := o.db.QueryRowContext(ctx, query, args...)

	rs := sicore.GetRowScanner(o.opts...)
	defer sicore.PutRowScanner(rs)

	err := rs.ScanPrimary(row, output)
	if err != nil {
		return err
	}

	return nil
}

func (o *SqlDB) QueryRowStruct(query string, output any, args ...any) error {
	return o.QueryRowContextStruct(context.Background(), query, output, args...)
}

func (o *SqlDB) QueryRowContextStruct(ctx context.Context, query string, output any, args ...any) error {
	rows, err := o.db.QueryContext(ctx, query, args...)
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
func (o *SqlDB) QueryStructs(query string, output any, args ...any) (int, error) {
	return o.QueryContextStructs(context.Background(), query, output, args...)
}

// QueryContextStructs queries a database with context then scan resultset into output of any type
func (o *SqlDB) QueryContextStructs(ctx context.Context, query string, output any, args ...any) (int, error) {
	rows, err := o.db.QueryContext(ctx, query, args...)
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

func (o *SqlDB) WithType(name string, typ sicore.SqlColType) *SqlDB {
	o.opts = append(o.opts, sicore.WithSqlColumnType(name, typ))
	return o
}

func (o *SqlDB) WithTypedBool(name string) *SqlDB {
	o.opts = append(o.opts, sicore.WithSqlColumnType(name, sicore.SqlColTypeBool))
	return o
}

func (o *SqlDB) appendRowScannerOpt(opt sicore.RowScannerOption) {
	o.opts = append(o.opts, opt)
}
