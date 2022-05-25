package sicore

import (
	"bufio"
	"io"
	"sync"
)

var (
	_readerPool = sync.Pool{
		New: func() interface{} {
			return new(bufio.Reader)
		},
	}

	_writerPool = sync.Pool{
		New: func() interface{} {
			return new(bufio.Writer)
		},
	}
)

func getBufioReader(r io.Reader) *bufio.Reader {
	br := _readerPool.Get().(*bufio.Reader)
	br.Reset(r)
	return br
}
func putBufioReader(br *bufio.Reader) {
	_readerPool.Put(br)
}

func getBufioWriter(w io.Writer) *bufio.Writer {
	bw := _writerPool.Get().(*bufio.Writer)
	bw.Reset(w)
	return bw
}
func putBufioWriter(bw *bufio.Writer) {
	_writerPool.Put(bw)
}
