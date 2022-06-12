package siwrap

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/go-wonk/si/sicore"
)

// HttpClient is a wrapper of http.Client
type HttpClient struct {
	client         *http.Client
	baseUrl        string
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
func (hc *HttpClient) DoReadBody(request *http.Request) ([]byte, int, error) {
	resp, err := hc.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	r := sicore.GetReader(resp.Body)
	defer sicore.PutReader(r)

	b, err := r.ReadAll()
	return b, resp.StatusCode, err
}

func (hc *HttpClient) request(medhod string, url string, header http.Header, body []byte) ([]byte, error) {
	var r *bytes.Reader
	var req *HttpRequest
	var err error
	if len(body) > 0 {
		r = sicore.GetBytesReader(body)
		defer sicore.PutBytesReader(r)
		req, err = GetRequest(medhod, url, r)
	} else {
		req, err = GetRequest(medhod, url, nil)
	}
	if err != nil {
		return nil, err
	}
	defer PutRequest(req)

	req.SetHeader(header)

	respBody, statusCode, err := hc.DoReadBody(req.Request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return respBody, fmt.Errorf("%d %s", statusCode, http.StatusText(statusCode))
	}

	return respBody, nil
}

func (hc *HttpClient) Request(medhod string, url string, header http.Header, body []byte) ([]byte, error) {
	return hc.request(medhod, hc.baseUrl+url, header, body)
}

func (hc *HttpClient) RequestGet(url string, header http.Header) ([]byte, error) {
	return hc.request(http.MethodGet, hc.baseUrl+url, header, nil)

}

func (hc *HttpClient) RequestPost(url string, header http.Header, body []byte) ([]byte, error) {
	return hc.request(http.MethodPost, hc.baseUrl+url, header, body)

}

func (hc *HttpClient) RequestPut(url string, header http.Header, body []byte) ([]byte, error) {
	return hc.request(http.MethodPut, hc.baseUrl+url, header, body)
}

func (hc *HttpClient) RequestDelete(url string, header http.Header, body []byte) ([]byte, error) {
	return hc.request(http.MethodDelete, hc.baseUrl+url, header, body)
}

func (hc *HttpClient) RequestHead(url string, header http.Header) ([]byte, error) {
	return hc.request(http.MethodHead, hc.baseUrl+url, header, nil)
}
