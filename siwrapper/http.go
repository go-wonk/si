package siwrapper

import (
	"net/http"

	"github.com/go-wonk/si/sicore"
)

const defaultBufferSize = 4096

type HttpClient struct {
	client         *http.Client
	defaultHeaders map[string]string
	bufferSize     int
}

func NewHttpClient(client *http.Client) *HttpClient {
	return NewHttpClientSizeWithHeader(client, defaultBufferSize, nil)
}
func NewHttpClientSize(client *http.Client, bufferSize int) *HttpClient {
	return NewHttpClientSizeWithHeader(client, bufferSize, nil)
}

func NewHttpClientWithHeader(client *http.Client, defaultHeaders map[string]string) *HttpClient {
	return NewHttpClientSizeWithHeader(client, defaultBufferSize, defaultHeaders)
}

func NewHttpClientSizeWithHeader(client *http.Client, bufferSize int, defaultHeaders map[string]string) *HttpClient {
	return &HttpClient{
		client:         client,
		defaultHeaders: defaultHeaders,
		bufferSize:     bufferSize,
	}
}

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

func (hc *HttpClient) Get(request *http.Request) ([]byte, error) {
	resp, err := hc.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := sicore.GetReaderSize(resp.Body, hc.bufferSize)
	b, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return b, nil
}
