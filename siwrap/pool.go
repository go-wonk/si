package siwrap

import (
	"database/sql"
	"sync"
)

var (
	_sqltxPool = sync.Pool{}
)

func getSqlTx(tx *sql.Tx) *SqlTx {
	g := _sqltxPool.Get()
	if g == nil {
		return newSqlTx(tx)
	}

	stx := g.(*SqlTx)
	stx.Reset(tx)
	return stx
}

func putSqlTx(sqlTx *SqlTx) {
	sqlTx.Reset(nil)
	_sqltxPool.Put(sqlTx)
}

func GetSqlTx(tx *sql.Tx) *SqlTx {
	return getSqlTx(tx)
}

func PutSqlTx(sqlTx *SqlTx) {
	putSqlTx(sqlTx)
}
