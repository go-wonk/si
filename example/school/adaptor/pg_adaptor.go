package adaptor

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-wonk/si/example/school/core"
	"github.com/go-wonk/si/sisql"
)

type pgStudentRepo struct {
	db *sisql.SqlDB
}

func NewPgStudentRepo(db *sql.DB) *pgStudentRepo {
	return &pgStudentRepo{db: sisql.NewSqlDB(db)}
}

func (o *pgStudentRepo) Add(student *core.Student, tx core.TxController) error {
	sqlTx := sisql.GetSqlTx(tx.(*sql.Tx))
	defer sisql.PutSqlTx(sqlTx)
	if sqlTx == nil {
		return errors.New("invalid tx")
	}

	insertQuery := `insert into student(email_address, name) values($1, $2)`
	_, err := sqlTx.Exec(insertQuery, student.EmailAddress, student.Name)
	if err != nil {
		return err
	}
	return nil
}

func (o *pgStudentRepo) Find(ID int) (*core.Student, error) {
	qry := `select * from student where id = $1`

	var output []core.Student
	n, err := o.db.QueryStructs(qry, &output, ID)
	if err != nil {
		return nil, err
	}
	if n == 1 {
		return &output[0], nil
	}

	return nil, nil
}

func (o *pgStudentRepo) FindAll() ([]core.Student, error) {
	qry := `select * from student order by id`

	var output []core.Student
	n, err := o.db.QueryStructs(qry, &output)
	if err != nil {
		return nil, err
	}

	return output[:n], nil
}

//
type pgBookRepo struct {
	db *sql.DB
}

func NewPgBookRepo(db *sql.DB) *pgBookRepo {
	return &pgBookRepo{db: db}
}

func (o *pgBookRepo) Add(book *core.Book, tx core.TxController) error {
	sqlTx := sisql.GetSqlTx(tx.(*sql.Tx))
	defer sisql.PutSqlTx(sqlTx)
	if sqlTx == nil {
		return errors.New("invalid tx")
	}

	insertQuery := `
		insert into book(name)
		values($1)
	`
	_, err := sqlTx.Exec(insertQuery, book.Name)
	if err != nil {
		return err
	}
	return nil
}

//
type pgBorrowingRepo struct {
	db *sql.DB
}

func NewPgBorrowingRepo(db *sql.DB) *pgBorrowingRepo {
	return &pgBorrowingRepo{db: db}
}

func (o *pgBorrowingRepo) Add(student *core.Student, book *core.Book, tx core.TxController) error {
	sqlTx := sisql.GetSqlTx(tx.(*sql.Tx))
	defer sisql.PutSqlTx(sqlTx)
	if sqlTx == nil {
		return errors.New("invalid tx")
	}

	insertQuery := `
		insert into borrowing(id, student_id, book_id, borrow_date)
		values($1, $2, $3, $4)
	`

	ID := core.GenerateID([]byte(fmt.Sprintf("%d_%d_%d", student.ID, book.ID, time.Now().UnixMilli())))
	_, err := sqlTx.Exec(insertQuery, ID, student.ID, book.ID, time.Now())
	if err != nil {
		return err
	}
	return nil
}
