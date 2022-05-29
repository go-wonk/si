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
