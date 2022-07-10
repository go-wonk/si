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
		Conn:            c,
		writeTimeout:    30 * time.Second,
		readTimeout:     30 * time.Second,
		writeBufferSize: 4096,
		readBufferSize:  4096,
		rw:              sicore.GetReadWriterWithReadWriter(c),
	}

	if err := conn.Reset(opts...); err != nil {
		return nil, err
	}

	// rw := sicore.GetReadWriterWithReadWriter(c)
	// rw.Reader.ApplyOptions(conn.readerOptions...)
	// rw.Writer.ApplyOptions(conn.writerOptions...)
	// conn.rw = rw

	return conn, nil
}

func (c *Conn) Reset(opts ...TcpOption) error {

	for _, o := range opts {
		if o == nil {
			continue
		}
		o.apply(c)
	}

	err := c.Conn.(*net.TCPConn).SetWriteBuffer(c.writeBufferSize)
	if err != nil {
		c.Conn.Close()
		return err
	}
	err = c.Conn.(*net.TCPConn).SetReadBuffer(c.readBufferSize)
	if err != nil {
		c.Conn.Close()
		return err
	}
	// c.Conn.(*net.TCPConn).SetKeepAlive(false)

	c.rw.Reader.Reset(c.Conn, c.readerOptions...)
	c.rw.Writer.Reset(c.Conn, c.writerOptions...)

	return nil
}

func (c *Conn) Close() error {
	sicore.PutReadWriter(c.rw)
	return c.Conn.Close()
}

func (c *Conn) Write(b []byte) (int, error) {
	err := c.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if err != nil {
		return 0, err
	}
	n, err := c.rw.Write(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (c *Conn) Read(b []byte) (int, error) {
	err := c.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return 0, err
	}
	n, err := c.rw.Read(b)
	if err != nil {
		return 0, err
	}
	return n, nil
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
