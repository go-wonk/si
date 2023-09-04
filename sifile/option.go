package sifile

import "github.com/go-wonk/si/v2/sicore"

// FileOption is an interface with apply method.
type FileOption interface {
	apply(f *File)
}

// FileOptionFunc wraps a function to conforms to FileOption interface
type FileOptionFunc func(f *File)

// apply implements FileOption's apply method.
func (o FileOptionFunc) apply(f *File) {
	o(f)
}

func WithReaderOpt(opt sicore.ReaderOption) FileOptionFunc {
	return FileOptionFunc(func(f *File) {
		f.appendReaderOpt(opt)
	})
}

func WithWriterOpt(opt sicore.WriterOption) FileOptionFunc {
	return FileOptionFunc(func(f *File) {
		f.appendWriterOpt(opt)
	})
}
