package main

import (
	"fmt"
	"log"

	"github.com/go-wonk/si/v2/sisql"
	_ "github.com/lib/pq"
)

func main() {

	type student struct {
		Name    string `si:"name"`
		Age     int    `si:"age"`
		Ignored string `si:"-"`
	}

	connStr := "host=testpghost port=5432 user=test password=test123 dbname=testdb sslmode=disable connect_timeout=60"
	driver := "postgres"
	db, err := sisql.Open(driver, connStr, sisql.WithTagKey("si"))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	query := `select 'wonk' as name, 20 as age`
	var s student
	err = db.QueryRowStruct(query, &s)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(s)

	var sl []student
	n, err := db.QueryStructs(query, &sl)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(n, sl)
}
