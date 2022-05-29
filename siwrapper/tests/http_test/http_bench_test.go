package http_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwrap"
	"github.com/stretchr/testify/assert"
)

func BenchmarkBasicClient_Get(b *testing.B) {
	if onlinetest != "1" {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	hc := siwrap.NewHttpClient(client)

	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("Content-type", "application/x-www-form-urlencoded")

		resp, err := hc.Do(request)
		siutils.AssertNilFailB(b, err)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			b.FailNow()
		}
		assert.EqualValues(b, "hello", string(body))
		resp.Body.Close()
	}
}

func BenchmarkHttpClient_Get(b *testing.B) {
	if onlinetest != "1" {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	hc := siwrap.NewHttpClient(client)

	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("Content-type", "application/x-www-form-urlencoded")

		body, err := hc.DoReadBody(request)
		siutils.AssertNilFailB(b, err)

		assert.EqualValues(b, "hello", string(body))
	}
}

func BenchmarkHttpClient_GetSize(b *testing.B) {
	if onlinetest != "1" {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	hc := siwrap.NewHttpClientSize(client, 512)

	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("Content-type", "application/x-www-form-urlencoded")

		body, err := hc.DoReadBody(request)
		siutils.AssertNilFailB(b, err)

		assert.EqualValues(b, "hello", string(body))
	}
}
