package si

import (
	"io"

	"github.com/go-wonk/si/sicore"
)

// ReadAll reads from r until EOF.
func ReadAll(r io.Reader) ([]byte, error) {
	return sicore.ReadAll(r)
}

// WriteAll writes input to w.
func WriteAll(w io.Writer, input []byte) (n int, err error) {
	return sicore.WriteAll(w, input)
}

// DecodeJson read src with json bytes then decode it into dst.
func DecodeJson(dst any, src io.Reader) error {
	sr := sicore.GetReader(src, sicore.SetJsonDecoder())
	defer sicore.PutReader(sr)
	return sr.Decode(dst)
}

// EncodeJson encode src into json bytes then write to dst.
func EncodeJson(dst io.Writer, src any) error {
	sw := sicore.GetWriter(dst, sicore.SetJsonEncoder())
	defer sicore.PutWriter(sw)
	return sw.EncodeFlush(src)
}
