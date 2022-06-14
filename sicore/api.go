package sicore

import "io"

// ReadAll reads from r until EOF.
func ReadAll(r io.Reader) ([]byte, error) {
	sr := GetReader(r)
	defer PutReader(sr)
	return sr.ReadAll()
}

// WriteAll writes b to w.
func WriteAll(w io.Writer, b []byte) (n int, err error) {
	sw := GetWriter(w)
	defer PutWriter(sw)

	return sw.WriteFlush(b)
}
