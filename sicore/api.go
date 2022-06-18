package sicore

import "io"

// ReadAll reads from r until EOF.
func ReadAll(r io.Reader) ([]byte, error) {
	sr := GetReader(r)
	defer PutReader(sr)
	return sr.ReadAll()
}

// WriteAll writes src to dst.
func WriteAll(dst io.Writer, src []byte) (n int, err error) {
	sw := GetWriter(dst)
	defer PutWriter(sw)

	return sw.WriteFlush(src)
}
