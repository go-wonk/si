package core

type StudentRepo interface {
	Add(student *Student, tx TxController) error
	Find(ID int) (*Student, error)
	FindAll() ([]Student, error)
}
type BookRepo interface {
	Add(book *Book, tx TxController) error
}
type BorrowingRepo interface {
	Add(student *Student, book *Book, tx TxController) error
}
