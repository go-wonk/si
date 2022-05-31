package sql_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var (
	onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))
	longtest, _   = strconv.ParseBool(os.Getenv("LONG_TEST"))
	// onlinetest, _ = strconv.ParseBool("1")

	db *sql.DB
)

type Table struct {
	Nil        string    `json:"nil" mapstructure:"nil"`
	Int2       int       `json:"int2_" mapstructure:"int2_"`
	Decimal    float64   `json:"decimal_" mapstructure:"decimal_"`
	Numeric    float64   `json:"numeric_" mapstructure:"numeric_"`
	Bigint     float64   `json:"bigint_" mapstructure:"bigint_"`
	CharArr    []uint8   `json:"char_arr_" mapstructure:"char_arr_"`
	VarcharArr []uint8   `json:"varchar_arr_" mapstructure:"varchar_arr_"`
	Bytea      []byte    `json:"bytea_" mapstructure:"bytea_"`
	Time       time.Time `json:"time_" mapstructure:"time_"`
}

func (t Table) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type TableList []Table

func (tl TableList) String() string {
	b, _ := json.Marshal(tl)
	return string(b)
}

func openDB() (*sql.DB, error) {
	connStr := "host=127.0.0.1 port=5432 user=test password=test123 dbname=testdb sslmode=disable connect_timeout=60"
	driver := "postgres"
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
