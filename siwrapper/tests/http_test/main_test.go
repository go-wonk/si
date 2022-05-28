package http_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var (
	onlinetest = os.Getenv("ONLINE_TEST")
	// onlinetest = "1"

	client *http.Client
)

func openClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       3,
		IdleConnTimeout:    time.Duration(15) * time.Second,
		DisableCompression: false,
	}

	client := &http.Client{
		Timeout:   time.Duration(15) * time.Second,
		Transport: tr,
	}
	return client
}

func setup() error {
	if onlinetest == "1" {
		client = openClient()
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
