package siwrap

import (
	"net/http"

	"github.com/go-wonk/si/sicore"
)

// HttpClient is a wrapper of http.Client
type HttpClient struct {
	client         *http.Client
	defaultHeaders map[string]string
}

// NewHttpClient returns default HttpClient
func NewHttpClient(client *http.Client) *HttpClient {
	return NewHttpClientWithHeader(client, nil)
}

// NewHttpClientWithHeader returns HttpClient with specified defaultHeaders that will be set on every request
func NewHttpClientWithHeader(client *http.Client, defaultHeaders map[string]string) *HttpClient {
	return &HttpClient{
		client:         client,
		defaultHeaders: defaultHeaders,
	}
}

// Do is a wrapper of http.Client.Do
func (hc *HttpClient) Do(request *http.Request) (*http.Response, error) {
	hc.setDefaultHeader(request)

	return hc.client.Do(request)
}

// setDefaultHeader sets defaultHeaders to request. It doesn't replace headers that are already assigned to `request`
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

	return r.ReadAll()
}

func (hc *HttpClient) PostReadBody(url string, header http.Header, body []byte) ([]byte, error) {
	r := sicore.GetBytesReader(body)
	defer sicore.PutBytesReader(r)

	// req, err := http.NewRequest(http.MethodPost, url, r)
	req, err := GetRequest(http.MethodPost, url, r)
	if err != nil {
		return nil, err
	}
	defer PutRequest(req)

	setHeader(req, header)

	return hc.DoReadBody(req)

}
