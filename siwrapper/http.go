package siwrapper

import (
	"io"
	"net/http"

	"github.com/go-wonk/si/sicore"
)

const defaultBufferSize = 4096

// HttpClient is a wrapper of http.Client
type HttpClient struct {
	client         *http.Client
	defaultHeaders map[string]string
	bufferSize     int
}

// NewHttpClient returns default HttpClient
func NewHttpClient(client *http.Client) *HttpClient {
	return NewHttpClientSizeWithHeader(client, defaultBufferSize, nil)
}

// NewHttpClientSize returns HttpClient with specified bufferSize
func NewHttpClientSize(client *http.Client, bufferSize int) *HttpClient {
	return NewHttpClientSizeWithHeader(client, bufferSize, nil)
}

// NewHttpClientWithHeader returns HttpClient with specified defaultHeaders that will be set on every request
func NewHttpClientWithHeader(client *http.Client, defaultHeaders map[string]string) *HttpClient {
	return NewHttpClientSizeWithHeader(client, defaultBufferSize, defaultHeaders)
}

// NewHttpClientSizeWithHeader returns HttpClient with specified bufferSize and defaultHeaders that will be set on every request
func NewHttpClientSizeWithHeader(client *http.Client, bufferSize int, defaultHeaders map[string]string) *HttpClient {
	return &HttpClient{
		client:         client,
		defaultHeaders: defaultHeaders,
		bufferSize:     bufferSize,
	}
}

// Do is a wrapper of http.Client.Do
func (hc *HttpClient) Do(request *http.Request) (*http.Response, error) {
	hc.setDefaultHeader(request)

	return hc.client.Do(request)
}

// setDefaultHeader sets defaultHeaders to request keeping already assigned headers of request
func (hc *HttpClient) setDefaultHeader(request *http.Request) {
	for k, v := range hc.defaultHeaders {
		if request.Header.Get(k) != "" {
			request.Header.Set(k, v)
		}
	}
}

// DoReadBody sends Do request and read all data from response.Body
func (hc *HttpClient) DoReadBody(request *http.Request) ([]byte, error) {
	resp, err := hc.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := sicore.GetReader(resp.Body)
	defer sicore.PutReader(r)
	b, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// func replaceByBufioReader(r io.Reader) io.Reader {
// 	if r != nil {
// 		var br *bufio.Reader
// 		if _, ok := r.(io.ByteReader); !ok {
// 			br = sicore.GetBufioReader(r)
// 			return br
// 		}
// 	}
// 	return r
// }

func NewGetRequest(url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodGet, url, body)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func NewPostRequest(url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return r, nil
}
