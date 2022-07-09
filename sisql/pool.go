package sisql

import (
	"database/sql"
	"sync"
)

var (
	_sqltxPool = sync.Pool{}
)

func getSqlTx(tx *sql.Tx, opts ...SqlTxOption) *SqlTx {
	g := _sqltxPool.Get()
	if g == nil {
		return newSqlTx(tx, opts...)
	}

	stx := g.(*SqlTx)
	stx.Reset(tx, opts...)
	return stx
}

func putSqlTx(sqlTx *SqlTx) {
	sqlTx.Reset(nil)
	_sqltxPool.Put(sqlTx)
}

func GetSqlTx(tx *sql.Tx, opts ...SqlTxOption) *SqlTx {
	return getSqlTx(tx, opts...)
}

func PutSqlTx(sqlTx *SqlTx) {
	putSqlTx(sqlTx)
}
