package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwrapper"
)

func TestHttpClient_Get(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.NotNilFail(t, client)

	hc := siwrapper.NewHttpClient(client)

	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
	siutils.NilFail(t, err)

	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	b, err := hc.Get(request)
	siutils.NilFail(t, err)

	fmt.Println(string(b))
}
