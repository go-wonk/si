package adaptor

import (
	"database/sql"

	"github.com/go-wonk/si/example/transaction/core"
)

type TxBeginner struct {
	db *sql.DB
}

func NewTxBeginner(db *sql.DB) *TxBeginner {
	return &TxBeginner{db}
}
func (t *TxBeginner) Begin() (core.TxController, error) {
	return t.db.Begin()
}
