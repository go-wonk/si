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
	return sicore.DecodeJson(dst, src)
}

// DecodeJsonCopied read src with json bytes then decode it into dst.
// It also write the data read from src into a bytes.Buffer then returns it.
func DecodeJsonCopied(dst any, src io.Reader) (*bytes.Buffer, error) {
	return sicore.DecodeJsonCopied(dst, src)
}

// EncodeJson encode src into json bytes then write to dst.
func EncodeJson(dst io.Writer, src any) error {
	return sicore.EncodeJson(dst, src)
}

// EncodeJsonCopied encode src into json bytes then write to dst.
// It also write encoded bytes of src to a bytes.Buffer then returns it.
func EncodeJsonCopied(dst io.Writer, src any) (*bytes.Buffer, error) {
	return sicore.EncodeJsonCopied(dst, src)
}
