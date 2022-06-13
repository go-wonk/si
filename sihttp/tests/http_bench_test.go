package sihttp_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/sihttp"
	"github.com/go-wonk/si/siutils"
	"github.com/stretchr/testify/assert"
)

func BenchmarkBasicClientGet(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	hc := sihttp.NewHttpClient(client)

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

func BenchmarkHttpClientGet(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	hc := sihttp.NewHttpClient(client)

	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("Content-type", "application/x-www-form-urlencoded")

		body, _, err := hc.DoRead(request)
		siutils.AssertNilFailB(b, err)

		assert.EqualValues(b, "hello", string(body))
	}
}

func BenchmarkHttpClientGetSize(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	hc := sihttp.NewHttpClient(client)

	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("Content-type", "application/x-www-form-urlencoded")

		body, _, err := hc.DoRead(request)
		siutils.AssertNilFailB(b, err)

		assert.EqualValues(b, "hello", string(body))
	}
}

const (
	testData = "a"
	// testDataRepeats = 1100
	testDataRepeats = 4096
	testUrl         = "https://127.0.0.1:8081/test/echo"
)

func BenchmarkReuseRequestPost(b *testing.B) {
	/*
		BenchmarkReuseRequestPost-8   	     192	  10939610 ns/op	    5433 B/op	      70 allocs/op
	*/
	if !onlinetest {
		b.Skip("skipping online tests")
	}

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)
		r := bytes.NewReader([]byte(sendData))

		req, err := http.NewRequest(http.MethodPost, url, r)
		if err != nil {
			b.FailNow()
		}

		req.Header.Set("custom_header", headerData)
		// req.URL.RawQuery = "bar=foo"

		resp, err := client.Do(req)
		if err != nil {
			b.FailNow()
		}

		io.ReadAll(resp.Body)
		// siutils.AssertNilFailB(b, err)
		// assert.EqualValues(b, sendData, string(respBody))
		// fmt.Println(string(respBody))

		resp.Body.Close()
	}
}

func BenchmarkReuseRequestPostWithRequestPool(b *testing.B) {
	/*
		BenchmarkReuseRequestPostWithRequestPool-8   	      79	  14926903 ns/op	    4773 B/op	      66 allocs/op
	*/
	if !onlinetest {
		b.Skip("skipping online tests")
	}

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)
		r := sicore.GetBytesReader([]byte(sendData))

		req, _ := sihttp.GetRequest(http.MethodPost, url, nil, r)
		// siutils.AssertNilFailB(b, err)

		req.Header.Set("custom_header", headerData)
		// req.URL.RawQuery = "bar=foo"

		resp, err := client.Do(req.Request)
		if err != nil {
			b.FailNow()
		}
		// siutils.AssertNilFailB(b, err)

		io.ReadAll(resp.Body)
		// siutils.AssertNilFailB(b, err)
		// assert.EqualValues(b, sendData, string(respBody))
		// fmt.Println(string(respBody))

		resp.Body.Close()
		sihttp.PutRequest(req)
		sicore.PutBytesReader(r)
	}
}

func BenchmarkBasicClientPost(b *testing.B) {
	// BenchmarkBasicClientPost-8   	    7180	    161840 ns/op	   32906 B/op	      88 allocs/op
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	hc := sihttp.NewHttpClient(client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(sendData)))
		if err != nil {
			b.FailNow()
		}

		request.Header.Set("custom_header", headerData)

		resp, err := hc.Do(request)
		if err != nil {
			b.FailNow()
		}

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			b.FailNow()
		}
		resp.Body.Close()
	}
}

func BenchmarkHttpClientPost(b *testing.B) {
	// BenchmarkReuseRequestPostWithRequestPool-8   	      79	  14926903 ns/op	    4773 B/op	      66 allocs/op
	// BenchmarkHttpClientPost-8   	    7731	    144920 ns/op	   19764 B/op	      79 allocs/op // no default encoding, just bytes
	// BenchmarkHttpClientPost2-8   	    7644	    145244 ns/op	   25158 B/op	      82 allocs/op // ReadWriter
	// BenchmarkHttpClientPost2-8   	    7677	    143666 ns/op	   20380 B/op	      82 allocs/op // Writer
	// BenchmarkHttpClientPost-8   	    7922	    146268 ns/op	   20371 B/op	      82 allocs/op
	if !onlinetest {
		b.Skip("skipping online tests")
	}

	client := sihttp.NewHttpClient(client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		header := make(http.Header)
		header["custom_header"] = []string{headerData}
		client.RequestPost(url, header, []byte(sendData))
	}
}

func BenchmarkHttpClientPost2(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}

	client := sihttp.NewHttpClient(client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		header := make(http.Header)
		header["custom_header"] = []string{headerData}
		client.RequestPost(url, header, []byte(sendData))
	}
}
