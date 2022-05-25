package main

import (
	"database/sql"
	"fmt"

	"github.com/go-wonk/si/example/transaction/adaptor"
	"github.com/go-wonk/si/example/transaction/core"
	_ "github.com/lib/pq"
)

var ()

func init() {

}

func main() {
	// sql storage
	connStr := "host=172.16.130.144 port=5432 user=test password=test123 dbname=testdb sslmode=disable connect_timeout=60"
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

	err = borrowingUsc.Borrow(&core.Student{ID: 1}, &core.Book{ID: 1})
	if err != nil {
		fmt.Println(err)
		return
	}

	bookUsc := core.NewBookUsecaseImpl(txBeginner, bookRepo)
	err = bookUsc.Add("Eva Armisen")
	if err != nil {
		fmt.Println(err)
		return
	}
}
