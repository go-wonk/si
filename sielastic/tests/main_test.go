package sielastic_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-wonk/si/v2/sielastic"
)

var (
	onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))
	// onlinetest, _ = strconv.ParseBool("1")
	// longtest, _ = strconv.ParseBool(os.Getenv("LONG_TEST"))
	// longtest, _ = strconv.ParseBool("1")
	client *elasticsearch.Client
)

func setup() error {
	var err error
	if onlinetest {
		client, err = sielastic.DefaultElasticsearchClient("http://testelastichost:9200", "daiso", "daisoasung")
		if err != nil {
			return err
		}
	}

	return nil
}

func shutdown() {
	// if client != nil {
	// 	client.shu
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
