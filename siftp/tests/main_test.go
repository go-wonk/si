package siftp_test

import (
	"fmt"
	"os"
	"testing"
)

var (
	onlinetest = false
	longtest   = false

// onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))
// longtest, _   = strconv.ParseBool(os.Getenv("LONG_TEST"))
)

func setup() error {
	return nil
}

func shutdown() {
	// do nothing yet
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
