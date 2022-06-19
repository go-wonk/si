package si

import (
	"bytes"
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

// DecodeJsonCopied read src with json bytes then decode it into dst.
// It also write the data read from src into a bytes.Buffer then returns it.
func DecodeJsonCopied(dst any, src io.Reader) (*bytes.Buffer, error) {
	bb := sicore.GetBytesBuffer(make([]byte, 0, 128))
	r := io.TeeReader(src, bb)
	sr := sicore.GetReader(r, sicore.SetJsonDecoder())
	defer sicore.PutReader(sr)
	return bb, sr.Decode(dst)
}

// EncodeJson encode src into json bytes then write to dst.
func EncodeJson(dst io.Writer, src any) error {
	sw := sicore.GetWriter(dst, sicore.SetJsonEncoder())
	defer sicore.PutWriter(sw)
	return sw.EncodeFlush(src)
}

// EncodeJsonCopied encode src into json bytes then write to dst.
// It also write encoded bytes of src to a bytes.Buffer then returns it.
func EncodeJsonCopied(dst io.Writer, src any) (*bytes.Buffer, error) {
	bb := sicore.GetBytesBuffer(make([]byte, 0, 128))
	mw := io.MultiWriter(dst, bb)
	sw := sicore.GetWriter(mw, sicore.SetJsonEncoder())
	defer sicore.PutWriter(sw)
	return bb, sw.EncodeFlush(src)
}
