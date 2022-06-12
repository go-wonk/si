package sihttp

import (
	"io"
	"sync"
)

// pool
var (
	_requestPool = sync.Pool{}
)

// GetRequest retrieves a request from a pool and returns it.
func GetRequest(method string, url string, body io.Reader) (*HttpRequest, error) {
	g := _requestPool.Get()
	if g == nil {
		return newHttpRequest(method, url, body)
	}
	req := g.(*HttpRequest)
	if err := req.Reset(method, url, body); err != nil {
		return nil, err
	}
	return req, nil

	// // Set Body, GetBody(), Content Length
	// req.setBody(body)

	// // Set method, url, host
	// if err := req.setMethodAndURL(method, url); err != nil {
	// 	return nil, err
	// }

	// // Clear headers
	// for k := range req.Header {
	// 	delete(req.Header, k)
	// }

	// // Clear trailers
	// for k := range req.Trailer {
	// 	delete(req.Trailer, k)
	// }

	// // Clear transfer encodings
	// req.TransferEncoding = req.TransferEncoding[:0]

	// return req, nil
}

func PutRequest(req *HttpRequest) {
	// req.Body.Close()
	req.Body = nil
	req.GetBody = nil
	_requestPool.Put(req)
}
