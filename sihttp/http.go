package sihttp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-wonk/si/sicore"
)

// HttpClient is a wrapper of http.Client
type HttpClient struct {
	client         *http.Client
	baseUrl        string
	defaultHeaders map[string]string

	requestOpts []RequestOption
	writerOpts  []sicore.WriterOption
	readerOpts  []sicore.ReaderOption
}

// NewHttpClient returns default HttpClient
func NewHttpClient(client *http.Client, opts ...RequestOption) *HttpClient {
	return NewHttpClientWithHeader(client, nil, opts...)
}

// NewHttpClientWithHeader returns HttpClient with specified defaultHeaders that will be set on every request
func NewHttpClientWithHeader(client *http.Client, defaultHeaders map[string]string, opts ...RequestOption) *HttpClient {
	c := &HttpClient{
		client:         client,
		defaultHeaders: defaultHeaders,
		requestOpts:    opts,
	}

	return c
}

func (hc *HttpClient) WithBaseUrl(baseUrl string) *HttpClient {
	hc.baseUrl = baseUrl
	return hc
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

func (hc *HttpClient) SetRequestOptions(opts ...RequestOption) {
	hc.requestOpts = opts
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

func (hc *HttpClient) request(method string, url string, header http.Header, body any) ([]byte, error) {
	var req *HttpRequest
	var err error

	req, err = GetRequest(method, url, header, body, hc.writerOpts...)
	if err != nil {
		return nil, err
	}
	defer PutRequest(req)

	for _, v := range hc.requestOpts {
		v.apply(req)
	}

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
	req, err = GetRequest(method, url, header, body, hc.writerOpts...)
	if err != nil {
		return err
	}
	defer PutRequest(req)

	for _, v := range hc.requestOpts {
		v.apply(req)
	}

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

func (hc *HttpClient) RequestPostReader(url string, header http.Header, body io.Reader) ([]byte, error) {
	return hc.request(http.MethodPost, hc.baseUrl+url, header, body)
}
func (hc *HttpClient) RequestPostDecodeReader(url string, header http.Header, body io.Reader, res any) error {
	return hc.requestDecode(http.MethodPost, hc.baseUrl+url, header, body, res)
}

func (hc *HttpClient) RequestPostFile(url string, header http.Header,
	params map[string]string, fileFieldName, fileName string) ([]byte, error) {

	// open file
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// create multipart.Writer
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile(fileFieldName, f.Name())
	if err != nil {
		return nil, err
	}

	// write file contents
	sr := sicore.GetReader(f)
	defer sicore.PutReader(sr)
	_, err = sr.WriteTo(w)
	if err != nil {
		return nil, err
	}

	// set Content-Type, overwrite existing Content-Type
	if header == nil {
		header = make(http.Header)
	}
	header["Content-Type"] = []string{mw.FormDataContentType()}

	// write params, this closes multipart.Writer
	for k, v := range params {
		mw.WriteField(k, v)
	}

	// close multipart writer
	if err = mw.Close(); err != nil {
		return nil, err
	}

	return hc.request(http.MethodPost, hc.baseUrl+url, header, buf)
}

// DefaultInsecureClient instantiate http.Client with InsecureSkipVerify set to true
func DefaultInsecureClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	return DefaultClient(tlsConfig)
}

// DefaultClient instantiate http.Client with input parameter `tlsConfig`
func DefaultClient(tlsConfig *tls.Config) *http.Client {

	dialer := &net.Dialer{Timeout: 3 * time.Second}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    time.Duration(15) * time.Second,
		DisableCompression: false,
		TLSClientConfig:    tlsConfig,
		DisableKeepAlives:  false,
		Dial:               dialer.Dial,
	}

	client := &http.Client{
		Timeout:   time.Duration(15) * time.Second,
		Transport: tr,
	}
	return client
}
