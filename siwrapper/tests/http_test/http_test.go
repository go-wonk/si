package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwrapper"
	"github.com/stretchr/testify/assert"
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

	b, err := hc.DoReadBody(request)
	siutils.NilFail(t, err)

	fmt.Println(string(b))
}

func TestNewGetRequest(t *testing.T) {
	// r, err := siwrapper.NewGetRequest("/test/hello", nil)
	// siutils.NilFail(t, err)

	hc := siwrapper.NewHttpClient(client)

	// hc.Do(r)

	pr, err := siwrapper.NewPostRequest("http://127.0.0.1:8080/test/echo", strings.NewReader("post request wrapper"))
	siutils.NilFail(t, err)

	body, err := hc.DoReadBody(pr)
	siutils.NilFail(t, err)

	assert.EqualValues(t, "post request wrapper", string(body))

	type Person struct {
		Name string `json:"name"`
		Age  uint8  `json:"age"`
	}

	var buf bytes.Buffer
	bw := sicore.GetBufioWriter(&buf)
	defer sicore.PutBufioWriter(bw)

	enc := json.NewEncoder(bw)
	siutils.NilFail(t, err)
	err = enc.Encode(&Person{"wonk", 20})
	siutils.NilFail(t, err)
	bw.Flush()
	fmt.Println(buf.String())
}
