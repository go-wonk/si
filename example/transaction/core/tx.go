package core

type TxController interface {
	Commit() error
	Rollback() error
}
type TxBeginner interface {
	Begin() (TxController, error)
}
