package sielastic_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-wonk/si/sielastic"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/tests/testmodels"
)

func TestElasticClient_IndexDocument(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, client)

	ec := sielastic.NewClient(client)
	// siutils.AssertNilFail(t, err)

	d := testmodels.Document{
		Name:      "my name is wonk",
		ID:        1040289,
		Timestamp: time.Now(),
	}
	body := d.String()
	res, err := ec.IndexDocument(context.Background(), "idx-test", "", []byte(body), true)
	siutils.AssertNilFail(t, err)

	fmt.Println(res)
}

func TestElasticClient_SearchDocuments(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, client)

	ec := sielastic.NewClient(client)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": "name",
			},
		},
	}
	res, err := ec.SearchDocuments(context.Background(), "idx-test", query)
	siutils.AssertNilFail(t, err)
	fmt.Println(res)
}
