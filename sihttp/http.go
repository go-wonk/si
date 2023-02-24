package sihttp

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/go-wonk/si/sicore"
)

// Client is a wrapper of http.Client
type Client struct {
	client         *http.Client
	baseUrl        string
	defaultHeaders map[string]string

	requestOpts []RequestOption
	writerOpts  []sicore.WriterOption
	readerOpts  []sicore.ReaderOption
}

// NewClient returns Client
func NewClient(client *http.Client, opts ...ClientOption) *Client {
	c := &Client{
		client: client,
	}
	for _, o := range opts {
		if o == nil {
			continue
		}
		o.apply(c)
	}

	return c
}

// Do is a wrapper of http.Client.Do
func (hc *Client) Do(request *http.Request) (*http.Response, error) {
	hc.setDefaultHeader(request)

	return hc.client.Do(request)
}

// setDefaultHeader sets defaultHeaders to request. It doesn't replace headers that are already assigned to `request`
func (hc *Client) setDefaultHeader(request *http.Request) {
	for k, v := range hc.defaultHeaders {
		if request.Header.Get(k) == "" {
			request.Header.Set(k, v)
		}
	}
}

func (hc *Client) appendRequestOption(opt RequestOption) {
	hc.requestOpts = append(hc.requestOpts, opt)
}
func (hc *Client) appendWriterOption(opt sicore.WriterOption) {
	hc.writerOpts = append(hc.writerOpts, opt)
}
func (hc *Client) appendReaderOption(opt sicore.ReaderOption) {
	hc.readerOpts = append(hc.readerOpts, opt)
}

// DoRead sends Do request and read all data from response.Body
func (hc *Client) DoRead(request *http.Request) ([]byte, int, error) {
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
func (hc *Client) DoDecode(request *http.Request, res any) (int, error) {
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

// setHeader sets `haeder` to underlying Request.
func setHeader(req *http.Request, header http.Header) {
	for k, val := range header {
		for i, v := range val {
			if i == 0 {
				req.Header.Set(k, v)
				continue
			}
			req.Header.Add(k, v)
		}
	}
}

func setQueries(req *http.Request, queries map[string]string) {
	if len(queries) > 0 {
		q := req.URL.Query()
		for k, v := range queries {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}

func (hc *Client) request(ctx context.Context, method string, url string,
	header http.Header, queries map[string]string, body any, opts ...RequestOption) ([]byte, error) {

	var req *http.Request
	var err error
	if r, ok := body.(io.Reader); ok {
		req, err = http.NewRequestWithContext(ctx, method, url, r)
	} else {
		w, buf := sicore.GetWriterAndBuffer(hc.writerOpts...)
		defer sicore.PutWriterAndBuffer(w, buf)
		if err := w.EncodeFlush(body); err != nil {
			return nil, err
		}
		req, err = http.NewRequestWithContext(ctx, method, url, buf)
	}
	if err != nil {
		return nil, err
	}

	setHeader(req, header)
	setQueries(req, queries)

	for _, v := range hc.requestOpts {
		v.apply(req)
	}

	for _, v := range opts {
		v.apply(req)
	}

	respBody, statusCode, err := hc.DoRead(req)
	if err != nil {
		return respBody, NewSiHttpError(statusCode, err.Error())
	}

	if statusCode >= 400 || statusCode < 100 {
		// TODO: should enable custom logic
		return respBody, NewSiHttpError(statusCode, http.StatusText(statusCode))
	}

	return respBody, nil
}

func (hc *Client) requestDecode(ctx context.Context, method string, url string, header http.Header, queries map[string]string, body any, res any, opts ...RequestOption) error {

	var req *http.Request
	var err error
	if r, ok := body.(io.Reader); ok {
		req, err = http.NewRequestWithContext(ctx, method, url, r)
	} else {
		w, buf := sicore.GetWriterAndBuffer(hc.writerOpts...)
		defer sicore.PutWriterAndBuffer(w, buf)
		if err := w.EncodeFlush(body); err != nil {
			return err
		}
		req, err = http.NewRequestWithContext(ctx, method, url, buf)
	}
	if err != nil {
		return err
	}

	setHeader(req, header)
	setQueries(req, queries)

	for _, v := range hc.requestOpts {
		v.apply(req)
	}

	for _, v := range opts {
		v.apply(req)
	}

	statusCode, err := hc.DoDecode(req, res)
	if err != nil {
		return NewSiHttpError(statusCode, err.Error())
	}

	if statusCode >= 400 || statusCode < 100 {
		// TODO: should enable custom logic
		return NewSiHttpError(statusCode, http.StatusText(statusCode))
	}

	return nil
}

func (hc *Client) Request(method string, url string, header http.Header, queries map[string]string, body []byte, opts ...RequestOption) ([]byte, error) {
	return hc.RequestContext(context.Background(), method, url, header, queries, body, opts...)
}
func (hc *Client) RequestContext(ctx context.Context, method string, url string, header http.Header, queries map[string]string, body []byte, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, method, hc.baseUrl+url, header, queries, body, opts...)
}

func (hc *Client) RequestGet(url string, header http.Header, queries map[string]string, opts ...RequestOption) ([]byte, error) {
	return hc.RequestGetContext(context.Background(), url, header, queries, opts...)
}
func (hc *Client) RequestGetContext(ctx context.Context, url string, header http.Header, queries map[string]string, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodGet, hc.baseUrl+url, header, queries, nil, opts...)
}

func (hc *Client) RequestPost(url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.RequestPostContext(context.Background(), url, header, body, opts...)
}
func (hc *Client) RequestPostContext(ctx context.Context, url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodPost, hc.baseUrl+url, header, nil, body, opts...)
}

func (hc *Client) RequestPut(url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.RequestPutContext(context.Background(), url, header, body, opts...)
}
func (hc *Client) RequestPutContext(ctx context.Context, url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodPut, hc.baseUrl+url, header, nil, body, opts...)
}

func (hc *Client) RequestDelete(url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.RequestDeleteContext(context.Background(), url, header, body, opts...)
}
func (hc *Client) RequestDeleteContext(ctx context.Context, url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodDelete, hc.baseUrl+url, header, nil, body, opts...)
}

func (hc *Client) RequestHead(url string, header http.Header, opts ...RequestOption) ([]byte, error) {
	return hc.RequestHeadContext(context.Background(), url, header, opts...)
}
func (hc *Client) RequestHeadContext(ctx context.Context, url string, header http.Header, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodHead, hc.baseUrl+url, header, nil, nil, opts...)
}

func (hc *Client) RequestPostReader(url string, header http.Header, body io.Reader, opts ...RequestOption) ([]byte, error) {
	return hc.RequestPostReaderContext(context.Background(), url, header, body, opts...)
}
func (hc *Client) RequestPostReaderContext(ctx context.Context, url string, header http.Header, body io.Reader, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodPost, hc.baseUrl+url, header, nil, body, opts...)
}

func (hc *Client) RequestDecode(method string, url string, header http.Header, queries map[string]string, body any, res any, opts ...RequestOption) error {
	return hc.RequestDecodeContext(context.Background(), http.MethodPost, url, header, queries, body, res, opts...)
}
func (hc *Client) RequestDecodeContext(ctx context.Context, method string, url string, header http.Header, queries map[string]string, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, method, hc.baseUrl+url, header, queries, body, res, opts...)
}
func (hc *Client) RequestGetDecode(url string, header http.Header, queries map[string]string, res any, opts ...RequestOption) error {
	return hc.RequestGetDecodeContext(context.Background(), url, header, queries, res, opts...)
}
func (hc *Client) RequestGetDecodeContext(ctx context.Context, url string, header http.Header, queries map[string]string, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodGet, hc.baseUrl+url, header, queries, nil, res, opts...)
}
func (hc *Client) RequestPostDecode(url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.RequestPostDecodeContext(context.Background(), url, header, body, res, opts...)
}
func (hc *Client) RequestPostDecodeContext(ctx context.Context, url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodPost, hc.baseUrl+url, header, nil, body, res, opts...)
}

func (hc *Client) RequestPutDecode(url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.RequestPutDecodeContext(context.Background(), url, header, body, res, opts...)
}
func (hc *Client) RequestPutDecodeContext(ctx context.Context, url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodPut, hc.baseUrl+url, header, nil, body, res, opts...)
}

func (hc *Client) RequestDeleteDecode(url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.RequestDeleteDecodeContext(context.Background(), url, header, body, res, opts...)
}
func (hc *Client) RequestDeleteDecodeContext(ctx context.Context, url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodDelete, hc.baseUrl+url, header, nil, body, res, opts...)
}

func (hc *Client) RequestHeadDecode(url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.RequestHeadDecodeContext(context.Background(), url, header, body, res, opts...)
}
func (hc *Client) RequestHeadDecodeContext(ctx context.Context, url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodHead, hc.baseUrl+url, header, nil, body, res, opts...)
}

func (hc *Client) RequestPostDecodeReader(url string, header http.Header, body io.Reader, res any, opts ...RequestOption) error {
	return hc.RequestPostDecodeReaderContext(context.Background(), url, header, body, res, opts...)
}
func (hc *Client) RequestPostDecodeReaderContext(ctx context.Context, url string, header http.Header, body io.Reader, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodPost, hc.baseUrl+url, header, nil, body, res, opts...)
}

func (hc *Client) RequestPostFile(url string, header http.Header,
	params map[string]string, fileFieldName, fileName string) ([]byte, error) {

	return hc.RequestPostFileContext(context.Background(), url, header, params, fileFieldName, fileName)
}

func (hc *Client) RequestPostFileContext(ctx context.Context, url string, header http.Header,
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

	return hc.request(ctx, http.MethodPost, hc.baseUrl+url, header, nil, buf)
}
