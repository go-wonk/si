package sicore

import (
	"bytes"
	"io"
)

// ReadAll reads from r until EOF.
func ReadAll(r io.Reader) ([]byte, error) {
	sr := GetReader(r)
	defer PutReader(sr)
	return sr.ReadAll()
}

// Deprecated
// WriteAll writes src to dst.
func WriteAll(dst io.Writer, src []byte) (n int, err error) {
	sw := GetWriter(dst)
	defer PutWriter(sw)

	return sw.WriteFlush(src)
}

// DecodeJson read src with json bytes then decode it into dst.
func DecodeJson(dst any, src io.Reader) error {
	sr := GetReader(src, SetJsonDecoder())
	defer PutReader(sr)
	return sr.Decode(dst)
}

// DecodeJsonCopied read src with json bytes then decode it into dst.
// It also write the data read from src into a bytes.Buffer then returns it.
func DecodeJsonCopied(dst any, src io.Reader) (*bytes.Buffer, error) {
	bb := GetBytesBuffer(nil)
	r := io.TeeReader(src, bb)
	sr := GetReader(r, SetJsonDecoder())
	defer PutReader(sr)
	return bb, sr.Decode(dst)
}

// EncodeJson encode src into json bytes then write to dst.
func EncodeJson(dst io.Writer, src any) error {
	sw := GetWriter(dst, SetJsonEncoder())
	defer PutWriter(sw)
	return sw.EncodeFlush(src)
}

// EncodeJsonCopied encode src into json bytes then write to dst.
// It also write encoded bytes of src to a bytes.Buffer then returns it.
func EncodeJsonCopied(dst io.Writer, src any) (*bytes.Buffer, error) {
	bb := GetBytesBuffer(nil)
	mw := io.MultiWriter(dst, bb)
	sw := GetWriter(mw, SetJsonEncoder())
	defer PutWriter(sw)
	return bb, sw.EncodeFlush(src)
}
