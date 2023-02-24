package sihttp

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/go-wonk/si/v2/sicore"
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

func (hc *Client) Request(method string, url string, header http.Header, queries map[string]string, body []byte, opts ...RequestOption) ([]byte, error) {
	return hc.RequestContext(context.Background(), method, url, header, queries, body, opts...)
}
func (hc *Client) RequestContext(ctx context.Context, method string, url string, header http.Header, queries map[string]string, body []byte, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, method, hc.baseUrl+url, header, queries, body, opts...)
}
func (hc *Client) RequestDecode(method string, url string, header http.Header, queries map[string]string, body any, res any, opts ...RequestOption) error {
	return hc.RequestDecodeContext(context.Background(), http.MethodPost, url, header, queries, body, res, opts...)
}
func (hc *Client) RequestDecodeContext(ctx context.Context, method string, url string, header http.Header, queries map[string]string, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, method, hc.baseUrl+url, header, queries, body, res, opts...)
}

func (hc *Client) Get(url string, header http.Header, queries map[string]string, opts ...RequestOption) ([]byte, error) {
	return hc.GetContext(context.Background(), url, header, queries, opts...)
}
func (hc *Client) GetContext(ctx context.Context, url string, header http.Header, queries map[string]string, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodGet, hc.baseUrl+url, header, queries, nil, opts...)
}
func (hc *Client) GetDecode(url string, header http.Header, queries map[string]string, res any, opts ...RequestOption) error {
	return hc.GetDecodeContext(context.Background(), url, header, queries, res, opts...)
}
func (hc *Client) GetDecodeContext(ctx context.Context, url string, header http.Header, queries map[string]string, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodGet, hc.baseUrl+url, header, queries, nil, res, opts...)
}

func (hc *Client) Post(url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.PostContext(context.Background(), url, header, body, opts...)
}
func (hc *Client) PostContext(ctx context.Context, url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodPost, hc.baseUrl+url, header, nil, body, opts...)
}
func (hc *Client) PostDecode(url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.PostDecodeContext(context.Background(), url, header, body, res, opts...)
}
func (hc *Client) PostDecodeContext(ctx context.Context, url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodPost, hc.baseUrl+url, header, nil, body, res, opts...)
}

func (hc *Client) Put(url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.PutContext(context.Background(), url, header, body, opts...)
}
func (hc *Client) PutContext(ctx context.Context, url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodPut, hc.baseUrl+url, header, nil, body, opts...)
}
func (hc *Client) PutDecode(url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.PutDecodeContext(context.Background(), url, header, body, res, opts...)
}
func (hc *Client) PutDecodeContext(ctx context.Context, url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodPut, hc.baseUrl+url, header, nil, body, res, opts...)
}

func (hc *Client) Delete(url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.DeleteContext(context.Background(), url, header, body, opts...)
}
func (hc *Client) DeleteContext(ctx context.Context, url string, header http.Header, body any, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodDelete, hc.baseUrl+url, header, nil, body, opts...)
}
func (hc *Client) DeleteDecode(url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.DeleteDecodeContext(context.Background(), url, header, body, res, opts...)
}
func (hc *Client) DeleteDecodeContext(ctx context.Context, url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodDelete, hc.baseUrl+url, header, nil, body, res, opts...)
}

func (hc *Client) Head(url string, header http.Header, opts ...RequestOption) ([]byte, error) {
	return hc.HeadContext(context.Background(), url, header, opts...)
}
func (hc *Client) HeadContext(ctx context.Context, url string, header http.Header, opts ...RequestOption) ([]byte, error) {
	return hc.request(ctx, http.MethodHead, hc.baseUrl+url, header, nil, nil, opts...)
}
func (hc *Client) HeadDecode(url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.HeadDecodeContext(context.Background(), url, header, body, res, opts...)
}
func (hc *Client) HeadDecodeContext(ctx context.Context, url string, header http.Header, body any, res any, opts ...RequestOption) error {
	return hc.requestDecode(ctx, http.MethodHead, hc.baseUrl+url, header, nil, body, res, opts...)
}

func (hc *Client) PostFile(url string, header http.Header,
	params map[string]string, fileFieldName, fileName string) ([]byte, error) {

	return hc.PostFileContext(context.Background(), url, header, params, fileFieldName, fileName)
}

func (hc *Client) PostFileContext(ctx context.Context, url string, header http.Header,
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
