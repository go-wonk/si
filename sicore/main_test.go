package sicore

import (
	"fmt"
	"os"
	"testing"
)

var (
	onlinetest = os.Getenv("ONLINE_TEST")
	// onlinetest = "1"

)

const (
	testDataFile = `{"name":"wonk","age":20,"email":"wonk@wonk.org"}` + "\n"
)

func setup() error {

	os.Mkdir("./tests/data", 0777)
	if onlinetest == "1" {

	}

	return nil
}

func shutdown() {
	// if db != nil {
	// 	db.Close()
	// }
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
