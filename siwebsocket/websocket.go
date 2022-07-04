package siwebsocket

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-wonk/si"
	"github.com/go-wonk/si/sicore"
	"github.com/gorilla/websocket"
)

func DefaultConn(u url.URL, header http.Header) (*websocket.Conn, *http.Response, error) {

	dialer := &websocket.Dialer{
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 3 * time.Second,
	}

	return dialer.Dial(u.String(), header)
}

func DefaultDialer(u url.URL, header http.Header) *websocket.Dialer {

	dialer := &websocket.Dialer{
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 3 * time.Second,
	}

	return dialer
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	// dialer *websocket.Dialer
	conn    *websocket.Conn
	handler MessageHandler

	// Time allowed to write a message to the peer.
	writeWait time.Duration
	// Time allowed to read the next pong message from the peer.
	readWait time.Duration
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod time.Duration
	// Maximum message size allowed from peer.
	maxMessageSize int

	// use ping/pong
	usePingPong bool

	data     chan []byte
	sendDone chan struct{}
	stopSend chan string

	readWg *sync.WaitGroup

	readErr  error
	writeErr error

	readOpts []sicore.ReaderOption
}

func (c *Client) appendReaderOpt(ro sicore.ReaderOption) {
	c.readOpts = append(c.readOpts, ro)
}

func (c *Client) ReadErr() error {
	return c.readErr
}

func (c *Client) WriteErr() error {
	return c.writeErr
}

func NewClientConfigured(conn *websocket.Conn, writeWait time.Duration, readWait time.Duration,
	maxMessageSize int, usePingPong bool, opts ...WebsocketOption) *Client {

	pingPeriod := (readWait * 9) / 10

	c := &Client{
		conn:    conn,
		handler: &NopMessageHandler{},

		writeWait:      writeWait,
		readWait:       readWait,
		pingPeriod:     pingPeriod,
		maxMessageSize: maxMessageSize,
		usePingPong:    usePingPong,

		data:     make(chan []byte),
		sendDone: make(chan struct{}),
		stopSend: make(chan string, 1),
		readWg:   &sync.WaitGroup{},
	}

	for _, o := range opts {
		o.apply(c)
	}

	go c.waitStopSend()
	c.readWg.Add(1)
	// go c.readPump()
	go c.writePump()

	return c
}

func NewClient(conn *websocket.Conn, opts ...WebsocketOption) *Client {
	writeWait := 10 * time.Second
	readWait := 60 * time.Second
	// pingPeriod := (pongWait * 9) / 10
	maxMessageSize := 512
	usePingPong := false

	return NewClientConfigured(conn, writeWait, readWait, maxMessageSize, usePingPong, opts...)
}

func (c *Client) SetMessageHandler(h MessageHandler) {
	c.handler = h
}

func (c *Client) waitStopSend() {
	<-c.stopSend
	close(c.sendDone)
}

var ErrStopChannelFull = errors.New("stop channel is full")

func (c *Client) Stop() error {
	select {
	case c.stopSend <- "stop":
	default:
		return ErrStopChannelFull
	}
	return nil
}

func (c *Client) Wait() {
	c.readWg.Wait()
}

var ErrDataChannelClosed = errors.New("send data channel closed")

func (c *Client) Send(b []byte) error {
	select {
	case <-c.sendDone:
		return ErrDataChannelClosed
	case c.data <- b:
	}
	return nil
}

func (c *Client) closeMessage(msg string) error {
	c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
	return c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, msg))
}

func (c *Client) ReadMessage() (messageType int, p []byte, err error) {
	var r io.Reader
	messageType, r, err = c.conn.NextReader()
	if err != nil {
		return messageType, nil, err
	}
	p, err = si.ReadAll(r)
	return messageType, p, err
}

func (c *Client) ReadPump() {
	defer func() {
		log.Println("return readPump")
		c.readWg.Done()
		c.conn.Close()
		c.Stop()
	}()

	c.conn.SetReadLimit(int64(c.maxMessageSize))
	c.conn.SetReadDeadline(time.Now().Add(c.readWait))
	if c.usePingPong {
		c.conn.SetPongHandler(func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(c.readWait))
			return nil
		})
	}

	cnt := 0 // TODO: for testing
	for {
		_, r, err := c.conn.NextReader()
		if err != nil {
			c.readErr = err
			return
		}
		c.handler.Handle(r)

		cnt++
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(c.pingPeriod)
	normalClose := false
	defer func() {
		// log.Println("return writePump")
		ticker.Stop()

		if !normalClose {
			// log.Println("abnormal write")
			c.conn.Close()
		}
		// c.conn.Close() // TODO: should be closed here?
		// Closing here causes losing messages sent by a peer. It is best to close the connection
		// on readPump so it read as much as possible.
	}()
	for {
		select {
		case <-c.sendDone:
			// Stop method has been called
			if err := c.closeMessage(""); err == nil {
				normalClose = true
			} else {
				// log.Println("write:", err)
				c.writeErr = err
			}
			return
		// nothing will close c.data channel
		// case message, ok := <-c.data:
		// 	if !ok {
		// 		// c.data channel closed.
		// 		c.closeMessage("")
		// 		return
		// 	}
		case message := <-c.data:
			c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				// log.Println(err)
				c.writeErr = err
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.data)
			for i := 0; i < n; i++ {
				w.Write(newline)
				msg := <-c.data
				w.Write(msg)
			}

			if err := w.Close(); err != nil {
				// log.Println(err)
				c.writeErr = err
				return
			}
		case <-ticker.C:
			if c.usePingPong {
				c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					c.writeErr = err
					return
				}
			}
		}
	}
}
