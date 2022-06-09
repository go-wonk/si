package sql_test

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

var (
	// onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))
	onlinetest, _ = strconv.ParseBool("1")
	longtest, _   = strconv.ParseBool(os.Getenv("LONG_TEST"))

	db *sql.DB
)

func openDB() (*sql.DB, error) {
	connStr := "host=192.168.0.92 port=5432 user=test password=test123 dbname=testdb sslmode=disable connect_timeout=60"
	driver := "postgres"
	// driver := "pgx"
	return sql.Open(driver, connStr)
}

func setup() error {
	if onlinetest {
		db, _ = openDB()
	}

	return nil
}

func shutdown() {
	if db != nil {
		db.Close()
	}

}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Println(err)
		shutdown()
		os.Exit(1)
	}

	exitCode := m.Run()

	shutdown()
	os.Exit(exitCode)
}
