package sitcp

import (
	"net"
	"time"

	"github.com/go-wonk/si/sicore"
)

func DialTimeout(addr string, timeout time.Duration, opts ...TcpOption) (*Conn, error) {
	c, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	conn, err := newConn(c, opts...)
	if err != nil {
		c.Close()
		return nil, err
	}
	return conn, nil
}

type Conn struct {
	net.Conn

	rw *sicore.ReadWriter

	writeTimeout    time.Duration
	readTimeout     time.Duration
	writeBufferSize int
	readBufferSize  int
	writerOptions   []sicore.WriterOption
	readerOptions   []sicore.ReaderOption
}

func newConn(c net.Conn, opts ...TcpOption) (*Conn, error) {
	conn := &Conn{
		writeTimeout:    30 * time.Second,
		readTimeout:     30 * time.Second,
		writeBufferSize: 4096,
		readBufferSize:  4096,
	}

	for _, o := range opts {
		if o == nil {
			continue
		}
		o.apply(conn)
	}

	err := c.(*net.TCPConn).SetWriteBuffer(conn.writeBufferSize)
	if err != nil {
		return nil, err
	}
	err = c.(*net.TCPConn).SetReadBuffer(conn.readBufferSize)
	if err != nil {
		return nil, err
	}

	rw := sicore.GetReadWriterWithReadWriter(c)
	rw.Reader.ApplyOptions(conn.readerOptions...)
	rw.Writer.ApplyOptions(conn.writerOptions...)

	conn.Conn = c
	conn.rw = rw

	return conn, nil
}

func (c *Conn) Close() error {
	sicore.PutReadWriter(c.rw)
	return c.Conn.Close()
}

func (c *Conn) Write(b []byte) (n int, err error) {
	err = c.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if err != nil {
		return
	}
	n, err = c.rw.Write(b)
	return
}

func (c *Conn) Read(b []byte) (n int, err error) {
	err = c.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return 0, err
	}
	n, err = c.rw.Read(b)
	return
}

func (c *Conn) appendReaderOption(opt sicore.ReaderOption) {
	c.readerOptions = append(c.readerOptions, opt)
}

func (c *Conn) appendWriterOption(opt sicore.WriterOption) {
	c.writerOptions = append(c.writerOptions, opt)
}

func (c *Conn) Request(b []byte) ([]byte, error) {
	err := c.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if err != nil {
		return nil, err
	}
	err = c.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return nil, err
	}

	return c.rw.Request(b)
}
