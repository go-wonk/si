package sihttp_test

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var (
	onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))
	longtest, _   = strconv.ParseBool(os.Getenv("LONG_TEST"))

	client *http.Client
)

func openClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	dialer := &net.Dialer{Timeout: 3 * time.Second}

	tr := &http.Transport{
		MaxIdleConns:       300,
		IdleConnTimeout:    time.Duration(15) * time.Second,
		DisableCompression: false,
		TLSClientConfig:    tlsConfig,
		DisableKeepAlives:  false,
		Dial:               dialer.Dial,
	}

	client := &http.Client{
		Timeout:   time.Duration(15) * time.Second,
		Transport: tr,
	}
	return client
}

func setup() error {
	if onlinetest {
		client = openClient()
	}

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
