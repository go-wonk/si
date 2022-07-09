package sisql

import "github.com/go-wonk/si/sicore"

// SqlOption is an interface with apply method.
type SqlOption interface {
	apply(db *SqlDB)
}

// SqlOptionFunc wraps a function to conforms to SqlOption interface.
type SqlOptionFunc func(db *SqlDB)

// apply implements SqlOption's apply method.
func (o SqlOptionFunc) apply(db *SqlDB) {
	o(db)
}

func WithRowScannerOpt(opt sicore.RowScannerOption) SqlOptionFunc {
	return SqlOptionFunc(func(db *SqlDB) {
		db.appendRowScannerOpt(opt)
	})
}

func WithTagKey(key string) SqlOptionFunc {
	return SqlOptionFunc(func(db *SqlDB) {
		db.appendRowScannerOpt(sicore.WithTagKey(key))
	})
}

func WithType(name string, typ sicore.SqlColType) SqlOptionFunc {
	return SqlOptionFunc(func(db *SqlDB) {
		db.appendRowScannerOpt(sicore.WithSqlColumnType(name, typ))
	})
}

// SqlTxOption is an interface with apply method.
type SqlTxOption interface {
	apply(db *SqlTx)
}

// SqlTxOptionFunc wraps a function to conforms to SqlTxOption interface.
type SqlTxOptionFunc func(db *SqlTx)

// apply implements SqlOption's apply method.
func (o SqlTxOptionFunc) apply(db *SqlTx) {
	o(db)
}

func WithTxRowScannerOpt(opt sicore.RowScannerOption) SqlTxOptionFunc {
	return SqlTxOptionFunc(func(db *SqlTx) {
		db.appendRowScannerOpt(opt)
	})
}

func WithTxTagKey(key string) SqlTxOptionFunc {
	return SqlTxOptionFunc(func(db *SqlTx) {
		db.appendRowScannerOpt(sicore.WithTagKey(key))
	})
}

func WithTxType(name string, typ sicore.SqlColType) SqlTxOptionFunc {
	return SqlTxOptionFunc(func(db *SqlTx) {
		db.appendRowScannerOpt(sicore.WithSqlColumnType(name, typ))
	})
}
