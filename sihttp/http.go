package sihttp

import (
	"bytes"
	"fmt"
	"io"
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
	c := &HttpClient{
		client:         client,
		defaultHeaders: defaultHeaders,
	}

	return c
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

func (hc *HttpClient) request(method string, url string, header http.Header, body io.Reader) ([]byte, error) {
	req, err := GetRequest(method, url, body)
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

func (hc *HttpClient) requestBytes(method string, url string, header http.Header, body []byte) ([]byte, error) {
	var r *bytes.Reader
	if len(body) > 0 {
		r = sicore.GetBytesReader(body)
		defer sicore.PutBytesReader(r)
		return hc.request(method, url, header, r)
	}

	// argument of body must be nil if len(body) == 0
	return hc.request(method, url, header, nil)

}

func (hc *HttpClient) Request(medhod string, url string, header http.Header, body []byte) ([]byte, error) {
	return hc.requestBytes(medhod, hc.baseUrl+url, header, body)
}

func (hc *HttpClient) RequestGet(url string, header http.Header) ([]byte, error) {
	return hc.requestBytes(http.MethodGet, hc.baseUrl+url, header, nil)

}

func (hc *HttpClient) RequestPost(url string, header http.Header, body []byte) ([]byte, error) {
	return hc.requestBytes(http.MethodPost, hc.baseUrl+url, header, body)

}

func (hc *HttpClient) RequestPut(url string, header http.Header, body []byte) ([]byte, error) {
	return hc.requestBytes(http.MethodPut, hc.baseUrl+url, header, body)
}

func (hc *HttpClient) RequestDelete(url string, header http.Header, body []byte) ([]byte, error) {
	return hc.requestBytes(http.MethodDelete, hc.baseUrl+url, header, body)
}

func (hc *HttpClient) RequestHead(url string, header http.Header) ([]byte, error) {
	return hc.requestBytes(http.MethodHead, hc.baseUrl+url, header, nil)
}

func (hc *HttpClient) RequestPostJson(url string, header http.Header, body any) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	s := sicore.GetWriter(buf, sicore.SetJsonEncoder())
	defer sicore.PutWriter(s)

	if err := s.EncodeFlush(body); err != nil {
		return nil, err
	}

	return hc.request(http.MethodPost, hc.baseUrl+url, header, buf)
}

func (hc *HttpClient) RequestPostJsonDecoded(url string, header http.Header, body any, out any) error {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	s := sicore.GetWriter(buf, sicore.SetJsonEncoder())
	defer sicore.PutWriter(s)

	if err := s.EncodeFlush(body); err != nil {
		return err
	}

	b, err := hc.request(http.MethodPost, hc.baseUrl+url, header, buf)
	if err != nil {
		return err
	}

	r := sicore.GetReader(bytes.NewReader(b), sicore.SetJsonDecoder())
	defer sicore.PutReader(r)

	if err = r.Decode(out); err != nil {
		return err
	}

	return nil
}
