package core

type StudentUsecase interface {
	Add(emailAddress, name string) error
	Find(ID int) (*Student, error)
	FindAll() ([]Student, error)
}
type StudentUsecaseImpl struct {
	txBeginner  TxBeginner
	studentRepo StudentRepo
}

func NewStudentUsecaseImpl(txBeginner TxBeginner, studentRepo StudentRepo) *StudentUsecaseImpl {
	u := &StudentUsecaseImpl{}
	u.txBeginner = txBeginner
	u.studentRepo = studentRepo
	return u
}

func (u *StudentUsecaseImpl) Add(emailAddress, name string) error {
	tx, err := u.txBeginner.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = u.studentRepo.Add(&Student{EmailAddress: emailAddress, Name: name}, tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *StudentUsecaseImpl) Find(ID int) (*Student, error) {
	return u.studentRepo.Find(ID)
}

func (u *StudentUsecaseImpl) FindAll() ([]Student, error) {
	return u.studentRepo.FindAll()
}

//
type BookUsecase interface {
	Add(name string) error
}
type BookUsecaseImpl struct {
	txBeginner TxBeginner
	bookRepo   BookRepo
}

func NewBookUsecaseImpl(txBeginner TxBeginner, bookRepo BookRepo) *BookUsecaseImpl {
	u := &BookUsecaseImpl{}
	u.txBeginner = txBeginner
	u.bookRepo = bookRepo
	return u
}

func (u *BookUsecaseImpl) Add(name string) error {
	tx, err := u.txBeginner.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = u.bookRepo.Add(&Book{Name: name}, tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

//
type BorrowingUsecase interface {
	Borrow(student *Student, book *Book) error
}
type BorrowingUsecaseImpl struct {
	txBeginner    TxBeginner
	borrowingRepo BorrowingRepo
	studentRepo   StudentRepo
	bookRepo      BookRepo
}

func NewBorrowingUsecaseImpl(txBeginner TxBeginner, borrowingRepo BorrowingRepo, studentRepo StudentRepo, bookRepo BookRepo) *BorrowingUsecaseImpl {
	u := &BorrowingUsecaseImpl{}
	u.txBeginner = txBeginner
	u.borrowingRepo = borrowingRepo
	u.studentRepo = studentRepo
	u.bookRepo = bookRepo
	return u
}

func (u *BorrowingUsecaseImpl) Borrow(student *Student, book *Book) error {
	tx, err := u.txBeginner.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = u.borrowingRepo.Add(student, book, tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
