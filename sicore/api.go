package sicore

import "io"

func ReadAll(r io.Reader) ([]byte, error) {
	sr := GetReader(r)
	defer PutReader(sr)
	return sr.ReadAll()
}

func WriteAll(w io.Writer, b []byte) (n int, err error) {
	sw := GetWriter(w)
	defer PutWriter(sw)

	return sw.WriteFlush(b)
}
