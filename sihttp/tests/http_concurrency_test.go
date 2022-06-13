package sihttp_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/go-wonk/si/sihttp"
	"github.com/stretchr/testify/assert"
)

// func requestt(wg *sync.WaitGroup, client *sihttp.HttpClient, method, url string, numRequests int, id int) error {
// 	// var err error
// 	wg.Add(1)
// 	go func(wg *sync.WaitGroup, id int) {
// 		for j := 0; j < numRequests; j++ {
// 			body := fmt.Sprintf("%s-%d-%d", "post hello", id, j)

// 			respBody, err := client.Request(method, url, nil, []byte(body))
// 			if err != nil {
// 				break
// 			}

// 			if !strings.EqualFold("hello", string(respBody)) {
// 				err = errors.New("unexpected response body")
// 				break
// 			}
// 		}
// 		wg.Done()
// 	}(wg, id)

// 	return nil
// }

func request(client *sihttp.HttpClient, method, url string, numRequests int, id int) error {
	var err error
	for j := 0; j < numRequests; j++ {

		var respBody []byte
		if method == http.MethodGet {
			respBody, err = client.Request(method, url, nil, nil)
			if err != nil {
				fmt.Println(err)
				break
			}
			if !strings.EqualFold("hello", string(respBody)) {
				fmt.Println(err)
				err = errors.New("unexpected response body")
				break
			}
		} else {
			body := fmt.Sprintf("%s-%d-%d", method+" hello", id, j)
			respBody, err = client.Request(method, url, nil, []byte(body))
			if err != nil {
				fmt.Println(err)
				break
			}
			var expected string
			if strings.HasSuffix(url, "/test/echo2") {
				expected = "2" + body
			} else if strings.HasSuffix(url, "/test/echo3") {
				expected = "3" + body
			} else if strings.HasSuffix(url, "/test/echo4") {
				expected = "4" + body
			} else {
				expected = body
			}
			if !strings.EqualFold(expected, string(respBody)) {
				err = errors.New("unexpected response body")
				fmt.Println(err)
				break
			}
		}

	}
	return err
}

func TestHttpClient_Concurrency_Request(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	if !longtest {
		t.Skip("skipping long tests")
	}

	client := sihttp.NewHttpClient(client)

	numRoutines := 20
	numRequests := 1000
	var wg sync.WaitGroup
	// wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {

		wg.Add(1)
		go func(wg *sync.WaitGroup, id int) {
			err := request(client, http.MethodGet, "http://127.0.0.1:8080/test/hello", numRequests, id)
			assert.Nil(t, err)
			wg.Done()
		}(&wg, i)

		wg.Add(1)
		go func(wg *sync.WaitGroup, id int) {
			err := request(client, http.MethodPost, "http://127.0.0.1:8080/test/echo", numRequests, id)
			assert.Nil(t, err)
			wg.Done()
		}(&wg, i)

		wg.Add(1)
		go func(wg *sync.WaitGroup, id int) {
			err := request(client, http.MethodPost, "http://127.0.0.1:8080/test/echo2", numRequests, id)
			assert.Nil(t, err)
			wg.Done()
		}(&wg, i)

		wg.Add(1)
		go func(wg *sync.WaitGroup, id int) {
			err := request(client, http.MethodPut, "http://127.0.0.1:8080/test/echo3", numRequests, id)
			assert.Nil(t, err)
			wg.Done()
		}(&wg, i)

		wg.Add(1)
		go func(wg *sync.WaitGroup, id int) {
			err := request(client, http.MethodDelete, "http://127.0.0.1:8080/test/echo4", numRequests, id)
			assert.Nil(t, err)
			wg.Done()
		}(&wg, i)

		// wg.Add(1)
		// go func(wg *sync.WaitGroup, i int) {
		// 	for j := 0; j < numRequests; j++ {
		// 		body := fmt.Sprintf("%s-%d-%d", "post hello", i, j)
		// 		url := "http://127.0.0.1:8080/test/hello"
		// 		respBody, err := client.RequestPost(url, nil, []byte(body))
		// 		siutils.AssertNilFail(t, err)

		// 		assert.EqualValues(t, "hello", string(respBody))
		// 	}
		// 	wg.Done()
		// 	fmt.Printf("done %d\n", i)
		// }(&wg, i)

		// wg.Add(1)
		// go func(wg *sync.WaitGroup, i int) {
		// 	for j := 0; j < numRequests; j++ {
		// 		body := fmt.Sprintf("%s-%d-%d", "post echo", i, j)
		// 		url := "http://127.0.0.1:8080/test/echo"
		// 		respBody, err := client.RequestPost(url, nil, []byte(body))
		// 		siutils.AssertNilFail(t, err)

		// 		assert.EqualValues(t, body, string(respBody))
		// 	}
		// 	wg.Done()
		// 	fmt.Printf("done %d\n", i)
		// }(&wg, i)

		// wg.Add(1)
		// go func(wg *sync.WaitGroup, i int) {
		// 	for j := 0; j < numRequests; j++ {
		// 		body := fmt.Sprintf("%s-%d-%d", "hello", i, j)
		// 		url := "http://127.0.0.1:8080/test/echo2"
		// 		respBody, err := client.RequestPost(url, nil, []byte(body))
		// 		siutils.AssertNilFail(t, err)

		// 		assert.EqualValues(t, "2"+body, string(respBody))
		// 	}
		// 	wg.Done()
		// 	fmt.Printf("done %d\n", i)
		// }(&wg, i)

		// wg.Add(1)
		// go func(wg *sync.WaitGroup, i int) {
		// 	for j := 0; j < numRequests; j++ {
		// 		body := fmt.Sprintf("%s-%d-%d", "hello", i, j)
		// 		url := "http://127.0.0.1:8080/test/echo3"
		// 		respBody, err := client.RequestPost(url, nil, []byte(body))
		// 		siutils.AssertNilFail(t, err)

		// 		assert.EqualValues(t, "3"+body, string(respBody))
		// 	}
		// 	wg.Done()
		// 	fmt.Printf("done %d\n", i)
		// }(&wg, i)
		// wg.Add(1)
		// go func(wg *sync.WaitGroup, i int) {
		// 	for j := 0; j < numRequests; j++ {
		// 		body := fmt.Sprintf("%s-%d-%d", "hello", i, j)
		// 		url := "http://127.0.0.1:8080/test/echo4"
		// 		respBody, err := client.RequestPost(url, nil, []byte(body))
		// 		siutils.AssertNilFail(t, err)

		// 		assert.EqualValues(t, "4"+body, string(respBody))
		// 	}
		// 	wg.Done()
		// 	fmt.Printf("done %d\n", i)
		// }(&wg, i)
	}
	wg.Wait()

}
