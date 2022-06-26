package sikafka_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

var (
	onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))
	// onlinetest, _ = strconv.ParseBool("1")
	// longtest, _ = strconv.ParseBool(os.Getenv("LONG_TEST"))
	// longtest, _ = strconv.ParseBool("1")

)

func setup() error {
	if onlinetest {
		// client = openClient()
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
