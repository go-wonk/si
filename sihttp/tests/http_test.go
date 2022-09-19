package sihttp_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/sihttp"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/tests/testmodels"
	"github.com/stretchr/testify/assert"
)

func TestHttpClientDo(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	siutils.AssertNotNilFail(t, client)

	hc := sihttp.NewClient(client)

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
	rw := sicore.GetReadWriterWithReadWriter(buf)
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
	rw := sicore.GetReadWriterWithReadWriter(buf)
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
	rw := sicore.GetReadWriterWithReadWriter(buf)
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
			rw := sicore.GetReadWriterWithReadWriter(buf)

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

func TestHttpClientRequestPostTls(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client)

	data := "hey"
	urls := []string{"http://127.0.0.1:8080/test/echo", "https://127.0.0.1:8081/test/echo"}
	for i := 0; i < 2; i++ {
		sendData := fmt.Sprintf("%s-%d", data, i)

		respBody, err := client.RequestPost(urls[i], nil, []byte(sendData))
		siutils.AssertNilFail(t, err)

		assert.EqualValues(t, sendData, string(respBody))
		fmt.Println(string(respBody))
	}
}

func TestHttpClientRequestGet(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client, sihttp.WithWriterOpt(sicore.SetJsonEncoder()))

	url := "http://127.0.0.1:8080/test/hello"

	queries := make(map[string]string)
	queries["name"] = "wonk"
	queries["kor"] = "길동"

	respBody, err := client.RequestGet(url, nil, queries)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, "hello", string(respBody))
	// fmt.Println(string(respBody))

}
func TestHttpClientRequestPost(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client)

	data := "hey"
	url := "http://127.0.0.1:8080/test/echo"

	sendData := fmt.Sprintf("%s-%d", data, 0)

	respBody, err := client.RequestPost(url, nil, []byte(sendData))
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))

}

func TestHttpClientRequestPut(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client)

	data := "hey"
	url := "http://127.0.0.1:8080/test/echo"

	sendData := fmt.Sprintf("%s-%d", data, 0)

	respBody, err := client.RequestPut(url, nil, []byte(sendData))
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))

}

func TestHttpClientRequestPostJsonDecoded(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client,
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()))

	url := "http://127.0.0.1:8080/test/echo"

	student := testmodels.Student{
		ID:           1,
		Name:         "wonk",
		EmailAddress: "wonk@wonk.org",
	}
	res := testmodels.Student{}
	err := client.RequestPostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)

	err = client.RequestPostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)
	// assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(res.String())

}

func TestHttpClientRequestPostFileData(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client)

	url := "http://127.0.0.1:8080/test/echo"

	f, err := os.OpenFile("./data/testfile.txt", os.O_RDONLY, 0777)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	header := make(http.Header)
	header["Content-Type"] = []string{"multipart/form-data"}

	res, err := client.RequestPostReader(url, header, f)
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res))

}

func TestHttpClientRequestPostReaderFile(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client)

	url := "http://127.0.0.1:8080/test/file/upload"

	f, err := os.OpenFile("./data/testfile.txt", os.O_RDONLY, 0777)
	siutils.AssertNilFail(t, err)
	defer f.Close()

	contents, err := io.ReadAll(f)
	siutils.AssertNilFail(t, err)

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	mw := multipart.NewWriter(buf)

	part, err := mw.CreateFormFile("file_to_upload", f.Name())
	siutils.AssertNilFail(t, err)
	part.Write(contents)

	mw.WriteField("nam", "wonk")

	header := make(http.Header)
	header["Content-Type"] = []string{mw.FormDataContentType()}

	err = mw.Close()
	siutils.AssertNilFail(t, err)

	// res, err := client.RequestPostFile(url, header, buf)
	res, err := client.RequestPostReader(url, header, buf)
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res))

}

func TestHttpClientRequestPostFile(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client,
		sihttp.WithRequestOpt(sihttp.WithHeaderHmac256("hmacKey", []byte("1234"))),
	)

	url := "http://127.0.0.1:8080/test/file/upload"

	res, err := client.RequestPostFile(url, nil, nil, "file_to_upload", "./data/testfile.txt")
	siutils.AssertNilFail(t, err)

	fmt.Println(string(res))

}

func TestHttpClientRequestGetWithHeaderHmac256(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client,
		sihttp.WithRequestOpt(sihttp.WithHeaderHmac256("hmac-hash", []byte("1234"))),
	)

	url := "http://127.0.0.1:8080/test/hello"

	respBody, err := client.RequestGet(url, nil, nil)
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, "hello", string(respBody))

}
func TestHttpClientRequestPostWithHeaderHmac256(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client,
		sihttp.WithRequestOpt(sihttp.WithHeaderHmac256("hmac-hash", []byte("1234"))),
	)

	data := "hey"
	url := "http://127.0.0.1:8080/test/echo"

	sendData := fmt.Sprintf("%s-%d", data, 0)

	respBody, err := client.RequestPost(url, nil, []byte(sendData))
	siutils.AssertNilFail(t, err)

	assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(string(respBody))

}

func TestHttpClientRequestPostJsonDecodedWithHeaderHmac256(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client,
		sihttp.WithRequestHeaderHmac256("hmacKey", []byte("1234")),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
	)
	// client.SetRequestOptions(sihttp.WithHeaderHmac256("hmacKey", []byte("1234")))
	// client.SetWriterOptions(sicore.SetJsonEncoder())
	// client.SetReaderOptions(sicore.SetJsonDecoder())

	url := "http://127.0.0.1:8080/test/echo"

	student := testmodels.Student{
		ID:           1,
		Name:         "wonk",
		EmailAddress: "wonk@wonk.org",
	}
	res := testmodels.Student{}
	err := client.RequestPostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)

	err = client.RequestPostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)
	// assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(res.String())

}

func TestHttpClientRequestPostJsonDecodedWithBearerToken(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client,
		sihttp.WithRequestOpt(sihttp.WithBearerToken("asdf")),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
	)

	url := "http://127.0.0.1:8080/test/echo"

	student := testmodels.Student{
		ID:           1,
		Name:         "wonk",
		EmailAddress: "wonk@wonk.org",
	}
	res := testmodels.Student{}
	err := client.RequestPostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)

	err = client.RequestPostDecode(url, nil, &student, &res)
	siutils.AssertNilFail(t, err)
	// assert.EqualValues(t, sendData, string(respBody))
	fmt.Println(res.String())

}

func TestWithBaseUrl(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	client := sihttp.NewClient(client,
		sihttp.WithBaseUrl("http://127.0.0.1:8080"),
	)

	url := "/test/echo"

	student := testmodels.Student{
		ID:           1,
		Name:         "wonk",
		EmailAddress: "wonk@wonk.org",
	}
	b, _ := json.Marshal(&student)
	res, err := client.RequestPost(url, nil, b)
	siutils.AssertNilFail(t, err)

	expected := `{"id":1,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false}`
	assert.EqualValues(t, expected, string(res))

}
