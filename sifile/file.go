package sifile

import (
	"os"

	"github.com/go-wonk/si/sicore"
)

// File is a wrapper of os.File
type File struct {
	*os.File

	r *sicore.Reader
	w *sicore.Writer
}

func NewFile(f *os.File) *File {
	r := sicore.GetReader(f)
	w := sicore.GetWriter(f)
	return &File{File: f, r: r, w: w}
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.r.Read(p)
}

func (f *File) ReadAll() ([]byte, error) {
	return f.r.ReadAll()
}

func (f *File) ReadAllFrom(offset int) ([]byte, error) {
	_, err := f.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return f.r.ReadAll()
}

func (f *File) Write(p []byte) (n int, err error) {
	return f.w.Write(p)
}

func (f *File) WriteFlush(p []byte) (n int, err error) {
	n, err = f.Write(p)
	if err != nil {
		return 0, err
	}

	if err = f.w.Flush(); err != nil {
		return 0, err
	}

	return n, err
}
