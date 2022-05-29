package http_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwrap"
	"github.com/stretchr/testify/assert"
)

func TestHttpClient_Get(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, client)

	hc := siwrap.NewHttpClient(client)

	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
	siutils.AssertNilFail(t, err)

	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	b, err := hc.DoReadBody(request)
	siutils.AssertNilFail(t, err)

	fmt.Println(string(b))
}

func TestNewGetRequest(t *testing.T) {
	if onlinetest != "1" {
		t.Skip("skipping online tests")
	}
	// r, err := siwrap.NewGetRequest("/test/hello", nil)
	// siutils.NilFail(t, err)

	hc := siwrap.NewHttpClient(client)

	// hc.Do(r)

	pr, err := siwrap.NewPostRequest("http://127.0.0.1:8080/test/echo", strings.NewReader("post request wrapper"))
	siutils.AssertNilFail(t, err)

	body, err := hc.DoReadBody(pr)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, "post request wrapper", string(body))
}

// func TestNewPostRequestJson(t *testing.T) {
// 	type Person struct {
// 		Name string `json:"name"`
// 		Age  uint8  `json:"age"`
// 	}

// 	hc := siwrap.NewHttpClient(client)

// 	pr, err := siwrap.NewPostRequestJson("http://127.0.0.1:8080/test/echo", &Person{"wonk", 20})
// 	siutils.NilFail(t, err)

// 	body, err := hc.DoReadBody(pr)
// 	siutils.NilFail(t, err)

// 	assert.EqualValues(t, `{"name":"wonk","age":20}`+"\n", string(body))

// }
