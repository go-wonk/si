package sihttp

import (
	"net/http"
	"sync"

	"github.com/go-wonk/si/sicore"
)

// pool
var (
	_requestPool = sync.Pool{}
)

// GetRequest retrieves a request from a pool and returns it.
func GetRequest(method string, url string, header http.Header, body any, opts ...sicore.WriterOption) (*HttpRequest, error) {
	g := _requestPool.Get()
	if g == nil {
		return newHttpRequest(method, url, header, body, opts...)
	}
	req := g.(*HttpRequest)
	if err := req.Reset(method, url, header, body, opts...); err != nil {
		return nil, err
	}
	return req, nil
}

func PutRequest(req *HttpRequest) {
	// req.Body.Close()
	// req.Body = nil
	// req.GetBody = nil
	_requestPool.Put(req)
}
