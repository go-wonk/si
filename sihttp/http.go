package sihttp

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

	writerOpts []sicore.WriterOption
	readerOpts []sicore.ReaderOption
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

func (hc *HttpClient) SetWriterOptions(opts ...sicore.WriterOption) {
	hc.writerOpts = opts
}
func (hc *HttpClient) SetReaderOptions(opts ...sicore.ReaderOption) {
	hc.readerOpts = opts
}

// DoRead sends Do request and read all data from response.Body
func (hc *HttpClient) DoRead(request *http.Request) ([]byte, int, error) {
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

// DoDecode sends Do request and decode response.Body
func (hc *HttpClient) DoDecode(request *http.Request, res any) (int, error) {
	resp, err := hc.Do(request)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	r := sicore.GetReader(resp.Body, hc.readerOpts...)
	defer sicore.PutReader(r)

	if err = r.Decode(res); err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}

// func (hc *HttpClient) request(method string, url string, header http.Header, body io.Reader) ([]byte, error) {
// 	req, err := GetRequest(method, url, body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer PutRequest(req)

// 	req.SetHeader(header)

// 	respBody, statusCode, err := hc.DoRead(req.Request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if statusCode != http.StatusOK {
// 		return respBody, fmt.Errorf("%d %s", statusCode, http.StatusText(statusCode))
// 	}

// 	return respBody, nil
// }

// func (hc *HttpClient) requestBytes(method string, url string, header http.Header, body []byte) ([]byte, error) {
// 	var r *bytes.Reader
// 	if len(body) > 0 {
// 		r = sicore.GetBytesReader(body)
// 		defer sicore.PutBytesReader(r)
// 		return hc.request(method, url, header, r)
// 	}

// 	// argument of body must be nil if len(body) == 0
// 	return hc.request(method, url, header, nil)
// }

// func (hc *HttpClient) requestEncode(method string, url string, header http.Header, body any) ([]byte, error) {
// 	buf := sicore.GetBytesBuffer(make([]byte, 0, 512))
// 	defer sicore.PutBytesBuffer(buf)

// 	w := sicore.GetWriter(buf, hc.writerOpts...)
// 	defer sicore.PutWriter(w)

// 	if err := w.EncodeFlush(body); err != nil {
// 		return nil, err
// 	}

// 	return hc.request(method, url, header, buf)
// }

func (hc *HttpClient) makeReqBody(buf *bytes.Buffer, body any) error {
	w := sicore.GetWriter(buf, hc.writerOpts...)
	defer sicore.PutWriter(w)

	if err := w.EncodeFlush(body); err != nil {
		return err
	}

	return nil
}

func (hc *HttpClient) request(method string, url string, header http.Header, body any) ([]byte, error) {
	var req *HttpRequest
	var err error
	if body == nil {
		req, err = GetRequest(method, url, nil)
	} else {
		buf := sicore.GetBytesBuffer(make([]byte, 0, 512))
		defer sicore.PutBytesBuffer(buf)

		if err := hc.makeReqBody(buf, body); err != nil {
			return nil, err
		}

		req, err = GetRequest(method, url, buf)
	}
	if err != nil {
		return nil, err
	}
	defer PutRequest(req)

	req.SetHeader(header)

	respBody, statusCode, err := hc.DoRead(req.Request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return respBody, fmt.Errorf("%d %s", statusCode, http.StatusText(statusCode))
	}

	return respBody, nil
}

func (hc *HttpClient) requestDecode(method string, url string, header http.Header, body any, res any) error {
	var req *HttpRequest
	var err error
	if body == nil {
		req, err = GetRequest(method, url, nil)
	} else {
		buf := sicore.GetBytesBuffer(make([]byte, 0, 512))
		defer sicore.PutBytesBuffer(buf)

		if err := hc.makeReqBody(buf, body); err != nil {
			return err
		}

		req, err = GetRequest(method, url, buf)
	}
	if err != nil {
		return err
	}
	defer PutRequest(req)

	req.SetHeader(header)

	statusCode, err := hc.DoDecode(req.Request, res)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("%d %s", statusCode, http.StatusText(statusCode))
	}

	return nil
}

func (hc *HttpClient) Request(method string, url string, header http.Header, body []byte) ([]byte, error) {
	return hc.request(method, hc.baseUrl+url, header, body)
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

func (hc *HttpClient) RequestDecode(method string, url string, header http.Header, body any, res any) error {
	return hc.requestDecode(http.MethodPost, hc.baseUrl+url, header, body, res)
}
func (hc *HttpClient) RequestGetDecode(url string, header http.Header, res any) error {
	return hc.requestDecode(http.MethodGet, hc.baseUrl+url, header, nil, res)
}
func (hc *HttpClient) RequestPostDecode(url string, header http.Header, body any, res any) error {
	return hc.requestDecode(http.MethodPost, hc.baseUrl+url, header, body, res)
}
