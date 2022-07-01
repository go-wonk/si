package sitcp

import (
	"net"
	"time"

	"github.com/go-wonk/si/sicore"
)

func DefaultTcpConn(addr string, dialTimeout, writeTimeout, readTimeout time.Duration, writeBuffer, readBuffer int) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, dialTimeout*time.Second)
	if err != nil {
		return nil, err
	}

	err = conn.SetWriteDeadline(time.Now().Add(writeTimeout * time.Second))
	if err != nil {
		return nil, err
	}
	err = conn.SetReadDeadline(time.Now().Add(readTimeout * time.Second))
	if err != nil {
		return nil, err
	}

	err = conn.(*net.TCPConn).SetWriteBuffer(writeBuffer)
	if err != nil {
		return nil, err
	}
	err = conn.(*net.TCPConn).SetReadBuffer(readBuffer)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type Conn struct {
	net.Conn
	writerOptions []sicore.WriterOption
	readerOptions []sicore.ReaderOption
	// addr                string
	// dialTimeoutSeconds  time.Duration
	// writeTimeoutSeconds time.Duration
	// readTimeoutSeconds  time.Duration
	// writeBufferSize     int
	// readBufferSize      int
}

func NewConn(conn net.Conn) *Conn {
	tcpConn := &Conn{Conn: conn}
	return tcpConn
}

func (c *Conn) Request(b []byte) ([]byte, error) {
	rw := sicore.GetReadWriterWithReadWriter(c)
	defer sicore.PutReadWriter(rw)
	rw.Reader.ApplyOptions(c.readerOptions...)
	rw.Writer.ApplyOptions(c.writerOptions...)

	return rw.Request(b)
}

func (c *Conn) SetWriterOption(opts ...sicore.WriterOption) {
	c.writerOptions = opts
}

func (c *Conn) SetReaderOption(opts ...sicore.ReaderOption) {
	c.readerOptions = opts
}
