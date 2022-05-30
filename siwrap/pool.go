package siwrap

import (
	"database/sql"
	"sync"

	"github.com/go-wonk/si/sicore"
)

var (
	_sqltxPool = sync.Pool{}
)

func getSqlTx(tx *sql.Tx, sc ...sicore.SqlColumn) *SqlTx {
	g := _sqltxPool.Get()
	if g == nil {
		return newSqlTx(tx, sc...)
	}

	stx := g.(*SqlTx)
	stx.Reset(tx, sc...)
	return stx
}

func putSqlTx(sqlTx *SqlTx) {
	sqlTx.Reset(nil)
	_sqltxPool.Put(sqlTx)
}

func GetSqlTx(tx *sql.Tx, sc ...sicore.SqlColumn) *SqlTx {
	return getSqlTx(tx, sc...)
}

func PutSqlTx(sqlTx *SqlTx) {
	putSqlTx(sqlTx)
}
