package http_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/siwrap"
	"github.com/stretchr/testify/assert"
)

func TestHttpClientDo(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, client)

	hc := siwrap.NewHttpClient(client)

	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/test/hello", nil)
	siutils.AssertNilFail(t, err)

	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(request)
	siutils.AssertNilFail(t, err)
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, "hello", string(b))
}

func TestCheckRequestState(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	data := "hey"

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	rw := sicore.GetReadWriterWithOptions(buf, nil, buf, nil)
	defer sicore.PutReadWriter(rw)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", rw)
	siutils.AssertNilFail(t, err)

	req.Header.Set("custom_header", "wonk")

	sendData := fmt.Sprintf("%s-%d", data, 0)
	rw.WriteFlush([]byte(sendData))
	resp, err := client.Do(req)
	siutils.AssertNilFail(t, err)

	respBody, err := io.ReadAll(resp.Body)
	siutils.AssertNilFail(t, err)
	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))
	resp.Body.Close()

	req2, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", rw)
	siutils.AssertNilFail(t, err)

	for k := range req.Header {
		delete(req.Header, k)
	}

	assert.EqualValues(t, req2, req)
}
func TestReuseRequest(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	data := "hey"

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	rw := sicore.GetReadWriterWithOptions(buf, nil, buf, nil)
	defer sicore.PutReadWriter(rw)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", rw)
	siutils.AssertNilFail(t, err)

	req.Header.Set("custom_header", "wonk")

	for i := 0; i < 10; i++ {
		sendData := fmt.Sprintf("%s-%d", data, i)
		rw.WriteFlush([]byte(sendData))
		resp, err := client.Do(req)
		siutils.AssertNilFail(t, err)

		respBody, err := io.ReadAll(resp.Body)
		siutils.AssertNilFail(t, err)
		assert.EqualValues(t, sendData, string(respBody))
		fmt.Println(string(respBody))

		resp.Body.Close()
	}

}

func TestReuseRequestInGoroutinePanic(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	t.Skip("skipping because this code panics")
	data := "hey"

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	rw := sicore.GetReadWriterWithOptions(buf, nil, buf, nil)
	defer sicore.PutReadWriter(rw)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", rw)
	siutils.AssertNilFail(t, err)

	var wg sync.WaitGroup
	for j := 0; j < 5; j++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			for i := 0; i < 10; i++ {
				sendData := fmt.Sprintf("%s-%d", data, i)

				req.Header.Set("custom_header", sendData)

				rw.WriteFlush([]byte(sendData))
				resp, err := client.Do(req)
				siutils.AssertNilFail(t, err)

				respBody, err := io.ReadAll(resp.Body)
				siutils.AssertNilFail(t, err)
				assert.EqualValues(t, sendData, string(respBody))
				fmt.Println(string(respBody))

				resp.Body.Close()
			}
			wg.Done()
		}(&wg)
	}
	wg.Wait()

}

func TestReuseRequestInGoroutine(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	data := "hey"

	var wg sync.WaitGroup
	for j := 0; j < 5; j++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, routineNumber int) {
			buf := bytes.NewBuffer(make([]byte, 0, 1024))
			rw := sicore.GetReadWriterWithOptions(buf, nil, buf, nil)

			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/test/echo", nil)
			siutils.AssertNilFail(t, err)

			req.Body = ioutil.NopCloser(rw)

			for i := 0; i < 10; i++ {
				sendData := fmt.Sprintf("%s-%d-%d", data, routineNumber, i)

				req.Header.Set("custom_header", sendData)

				rw.WriteFlush([]byte(sendData))
				resp, err := client.Do(req)
				siutils.AssertNilFail(t, err)

				respBody, err := io.ReadAll(resp.Body)
				siutils.AssertNilFail(t, err)
				assert.EqualValues(t, sendData, string(respBody))
				fmt.Println(string(respBody))

				resp.Body.Close()
			}

			sicore.PutReadWriter(rw)
			wg.Done()
		}(&wg, j)
	}
	wg.Wait()

}

func TestReuseRequestWithRequestPool(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	data := "hey"

	rw := bytes.NewBuffer(make([]byte, 0, 1024))

	urls := []string{"http://127.0.0.1:8080/test/echo", "https://127.0.0.1:8081/test/echo"}
	for i := 0; i < 2; i++ {
		sendData := fmt.Sprintf("%s-%d", data, i)
		rw.Write([]byte(sendData))

		req, err := siwrap.GetRequest(http.MethodPost, urls[i], rw)
		siutils.AssertNilFail(t, err)

		//////////////////////////////////////////////////////////
		// Check if pooled request is porperly reset
		expected, err := http.NewRequest(http.MethodPost, urls[i], rw)
		siutils.AssertNilFail(t, err)
		assert.EqualValues(t, expected.Method, req.Method)
		assert.EqualValues(t, expected.URL, req.URL)
		assert.EqualValues(t, expected.Proto, req.Proto)
		assert.EqualValues(t, expected.ProtoMajor, req.ProtoMajor)
		assert.EqualValues(t, expected.ProtoMinor, req.ProtoMinor)
		assert.EqualValues(t, expected.Header, req.Header)
		assert.EqualValues(t, expected.Body, req.Body)
		assert.EqualValues(t, expected.Host, req.Host)
		assert.EqualValues(t, expected.ContentLength, req.ContentLength)

		assert.EqualValues(t, expected.TransferEncoding, req.TransferEncoding)
		assert.EqualValues(t, expected.Trailer, req.Trailer)             // For client, once the body returns EOF(read all), the caller must not mutate Trailer.
		assert.EqualValues(t, expected.Close, req.Close)                 // DO NOT SET THIS ON CLIENT
		assert.EqualValues(t, expected.Form, req.Form)                   // DO NOT SET THIS ON CLIENT
		assert.EqualValues(t, expected.PostForm, req.PostForm)           // DO NOT SET THIS ON CLIENT
		assert.EqualValues(t, expected.MultipartForm, req.MultipartForm) // DO NOT SET THIS ON CLIENT
		assert.EqualValues(t, expected.RemoteAddr, req.RemoteAddr)       // DO NOT SET THIS ON CLIENT
		assert.EqualValues(t, expected.RequestURI, req.RequestURI)       // DO NOT SET THIS ON CLIENT
		//////////////////////////////////////////////////////////

		req.Header.Set("custom_header", sendData)
		req.URL.RawQuery = "bar=foo"

		resp, err := client.Do(req)
		siutils.AssertNilFail(t, err)

		respBody, err := io.ReadAll(resp.Body)
		siutils.AssertNilFail(t, err)
		assert.EqualValues(t, sendData, string(respBody))
		fmt.Println(string(respBody))

		resp.Body.Close()

		siwrap.PutRequest(req)
	}
}

func TestHttpClientPostReadBody(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := siwrap.NewHttpClient(client)

	data := "hey"
	urls := []string{"http://127.0.0.1:8080/test/echo", "https://127.0.0.1:8081/test/echo"}
	for i := 0; i < 2; i++ {
		sendData := fmt.Sprintf("%s-%d", data, i)

		respBody, err := client.PostReadBody(urls[i], nil, []byte(sendData))
		siutils.AssertNilFail(t, err)

		assert.EqualValues(t, sendData, string(respBody))
		fmt.Println(string(respBody))
	}
}
