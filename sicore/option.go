package sicore

type WriterOption interface {
	apply(w *Writer)
}

type WriterOptionFunc func(*Writer)

func (s WriterOptionFunc) apply(w *Writer) {
	s(w)
}

type ReaderOption interface {
	apply(r *Reader)
}
type ReaderOptionFunc func(*Reader)

func (o ReaderOptionFunc) apply(r *Reader) {
	o(r)
}

type RowScannerOption interface {
	apply(rs *RowScanner)
}
type RowScannerOptionFunc func(rs *RowScanner)

func (o RowScannerOptionFunc) apply(rs *RowScanner) {
	o(rs)
}

// WithTagKey sets RowScanner's tagKey
func WithTagKey(key string) RowScannerOption {
	return RowScannerOptionFunc(func(rs *RowScanner) {
		rs.SetTagKey(key)
	})
}
