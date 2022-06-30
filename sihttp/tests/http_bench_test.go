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

// BenchmarkBasicClientGet-8   	     794	   1445329 ns/op	    4925 B/op	      57 allocs/op
func BenchmarkHttpClient_DefaultGet(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("Content-type", "application/x-www-form-urlencoded")

		resp, err := client.Do(request)
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

// BenchmarkHttpClientGet-8   	     939	   1463128 ns/op	    4420 B/op	      57 allocs/op
func BenchmarkHttpClient_DoRead(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	hc := sihttp.NewHttpClient(client)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("Content-type", "application/x-www-form-urlencoded")

		body, _, err := hc.DoRead(request)
		siutils.AssertNilFailB(b, err)

		assert.EqualValues(b, "hello", string(body))
	}
}

// BenchmarkHttpClient_RequestGet-8   	     825	   1494142 ns/op	    3777 B/op	      55 allocs/op
func BenchmarkHttpClient_RequestGet(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	hc := sihttp.NewHttpClient(client)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		header := make(http.Header)
		header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
		body, err := hc.RequestGet("http://127.0.0.1:8080/test/hello", header, nil)
		siutils.AssertNilFailB(b, err)

		assert.EqualValues(b, "hello", string(body))
	}
}

const (
	testData = "a"
	// testDataRepeats = 128
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

// BenchmarkHttpClient_DefaultPost-8   	     342	   4046098 ns/op	   33199 B/op	      88 allocs/op
// BenchmarkHttpClient_DefaultPost-8   	    1555	    759026 ns/op	    5857 B/op	      74 allocs/op
func BenchmarkHttpClient_DefaultPost(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(sendData)))
		siutils.AssertNilFailB(b, err)

		request.Header.Set("custom_header", headerData)

		resp, err := client.Do(request)
		siutils.AssertNilFailB(b, err)

		_, err = io.ReadAll(resp.Body)
		siutils.AssertNilFailB(b, err)

		resp.Body.Close()
	}
}

// BenchmarkHttpClient_DefaultPost_WithPool-8   	     370	   3142342 ns/op	   33095 B/op	      87 allocs/op
// BenchmarkHttpClient_DefaultPost_WithPool-8   	    1348	    884274 ns/op	    5878 B/op	      74 allocs/op
func BenchmarkHttpClient_DefaultPost_WithPool(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		buf := sicore.GetBytesReader([]byte(sendData))
		request, err := http.NewRequest(http.MethodPost, url, buf)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("custom_header", headerData)

		resp, err := client.Do(request)
		siutils.AssertNilFailB(b, err)

		_, err = io.ReadAll(resp.Body)
		siutils.AssertNilFailB(b, err)

		resp.Body.Close()
		sicore.PutBytesReader(buf)
	}
}

// BenchmarkHttpClient_DefaultPost_WithPoolAndDoRead-8   	     456	   2813064 ns/op	   20599 B/op	      81 allocs/op
// BenchmarkHttpClient_DefaultPost_WithPoolAndDoRead-8   	    1231	    981608 ns/op	    5478 B/op	      73 allocs/op
func BenchmarkHttpClient_DefaultPost_WithPoolAndDoRead(b *testing.B) {
	if !onlinetest {
		b.Skip("skipping online tests")
	}
	siutils.AssertNotNilFailB(b, client)

	data := strings.Repeat(testData, testDataRepeats)
	url := testUrl

	client := sihttp.NewHttpClient(client)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		sendData := fmt.Sprintf("%s-%d", data, i)
		headerData := fmt.Sprintf("%d", i)

		buf := sicore.GetBytesReader([]byte(sendData))
		request, err := http.NewRequest(http.MethodPost, url, buf)
		siutils.AssertNilFailB(b, err)

		request.Header.Set("custom_header", headerData)

		_, _, err = client.DoRead(request)
		siutils.AssertNilFailB(b, err)

		sicore.PutBytesReader(buf)
	}
}

// BenchmarkHttpClient_RequestPost-8   	     400	   3210657 ns/op	   20221 B/op	      80 allocs/op
// BenchmarkHttpClient_RequestPost-8   	    1384	    733201 ns/op	    4858 B/op	      72 allocs/op
func BenchmarkHttpClient_RequestPost(b *testing.B) {
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
		_, err := client.RequestPost(url, header, []byte(sendData))
		siutils.AssertNilFailB(b, err)
	}
}
