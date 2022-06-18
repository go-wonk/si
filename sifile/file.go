package sifile

import (
	"io"
	"os"

	"github.com/go-wonk/si/sicore"
)

// File is a wrapper of os.File
type File struct {
	*os.File

	rw *sicore.ReadWriter
}

func OpenFile(name string, flag int, perm os.FileMode) (*File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

func newFile(f *os.File) *File {
	rw := sicore.GetReadWriterWithReadWriter(f)
	// w := sicore.GetWriter(f)
	return &File{File: f, rw: rw}
}

func (f *File) Chdir() error {
	return f.File.Chdir()
}

func (f *File) Chmod(mode os.FileMode) error {
	return f.File.Chmod(mode)
}

func (f *File) Chown(uid, gid int) error {
	return f.File.Chown(uid, gid)
}

func (f *File) Close() error {
	sicore.PutReadWriter(f.rw)
	return f.File.Close()
}

func (f *File) Fd() uintptr {
	return f.File.Fd()
}

func (f *File) Name() string {
	return f.File.Name()
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.rw.Read(p)
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	return f.File.ReadAt(b, off)
}

func (f *File) ReadDir(n int) ([]os.DirEntry, error) {
	return f.File.ReadDir(n)
}

func (f *File) ReadFrom(r io.Reader) (n int64, err error) {
	return f.rw.ReadFrom(r)
}

func (f *File) Readdir(n int) ([]os.FileInfo, error) {
	return f.File.Readdir(n)
}

func (f *File) Readdirnames(n int) (names []string, err error) {
	return f.File.Readdirnames(n)
}

func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	return f.File.Seek(offset, whence)
}

func (f *File) Write(b []byte) (n int, err error) {
	return f.rw.Write(b)
}

func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	return f.File.WriteAt(b, off)
}

func (f *File) WriteString(s string) (n int, err error) {
	return f.rw.WriteString(s)
}

func (f *File) ReadAll() ([]byte, error) {
	return f.rw.ReadAll()
}

func (f *File) ReadAllFrom(offset int64) ([]byte, error) {
	_, err := f.Seek(offset, 0)
	if err != nil {
		return nil, err
	}
	return f.rw.ReadAll()
}

func (f *File) WriteFlush(p []byte) (n int, err error) {
	n, err = f.Write(p)
	if err != nil {
		return 0, err
	}

	if err = f.rw.Flush(); err != nil {
		return 0, err
	}

	return n, err
}
