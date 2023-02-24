package main

import (
	"database/sql"
	"fmt"

	"github.com/go-wonk/si/v2/example/school/adaptor"
	"github.com/go-wonk/si/v2/example/school/core"
	_ "github.com/lib/pq"
)

var ()

func init() {

}

func main() {
	// sql storage
	connStr := "host=testpghost port=5432 user=test password=test123 dbname=testdb sslmode=disable connect_timeout=60"
	driver := "postgres"
	db, err := sql.Open(driver, connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	txBeginner := adaptor.NewTxBeginner(db)
	studentRepo := adaptor.NewPgStudentRepo(db)
	bookRepo := adaptor.NewPgBookRepo(db)
	borrowingRepo := adaptor.NewPgBorrowingRepo(db)

	borrowingUsc := core.NewBorrowingUsecaseImpl(txBeginner, borrowingRepo, studentRepo, bookRepo)
	bookUsc := core.NewBookUsecaseImpl(txBeginner, bookRepo)

	err = borrowingUsc.Borrow(&core.Student{ID: 1}, &core.Book{ID: 1})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bookUsc.Add("Eva Armisen")
	if err != nil {
		fmt.Println(err)
		return
	}
}
