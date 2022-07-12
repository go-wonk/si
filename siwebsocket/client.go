package siwebsocket

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-wonk/si"
	"github.com/go-wonk/si/sicore"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
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

// Client starts two go-routines to write and read messages upon creation with
// NewClientCongfigured or NewClient function. SendMessage method is used to
// write messages to the server, and read messages are handled by underlying
// handler(MessageHandler). Currently, it is not designed to work with request
// (send and receive) pattern. Sending and reading messages are handled separately,
// you cannot verify if a request is correctly handled by the other end point's response.
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
	msg      chan *msg
	sendDone chan struct{}
	stopSend chan string

	readWg *sync.WaitGroup

	readErr  error
	writeErr error

	readOpts []sicore.ReaderOption

	id string

	// for server only
	hub         sicore.Hub
	userID      string
	userGroupID string
}

func NewClient(conn *websocket.Conn, opts ...ClientOption) (*Client, error) {

	defaultReadWait := 60 * time.Second
	c := &Client{
		conn:    conn,
		handler: &NopMessageHandler{},

		writeWait:      10 * time.Second,
		readWait:       defaultReadWait,
		pingPeriod:     (defaultReadWait * 9) / 10,
		maxMessageSize: 4096,
		usePingPong:    false,

		data:     make(chan []byte),
		msg:      make(chan *msg),
		sendDone: make(chan struct{}),
		stopSend: make(chan string, 1),
		readWg:   &sync.WaitGroup{},

		hub: &sicore.NopHub{},
	}

	for _, o := range opts {
		o.apply(c)
	}

	if c.id == "" {
		c.id = uuid.New().String()
	}

	go c.waitStopSend()
	go c.writePump()
	c.readWg.Add(1)
	go c.readPump()

	err := c.hub.Add(c)
	if err != nil {
		c.Stop()
		c.Wait()
		return nil, err
	}

	return c, nil
}

type upgraderConfig struct {
	handshakeTimeout                time.Duration
	readBufferSize, writeBufferSize int
	writeBufferPool                 websocket.BufferPool
	subprotocols                    []string
	errorFunc                       func(w http.ResponseWriter, r *http.Request, status int, reason error)
	checkOrigin                     func(r *http.Request) bool
	enableCompression               bool
}

func GetUpgradeConfig(opts ...UpgraderOption) *upgraderConfig {
	u := &upgraderConfig{}
	for _, o := range opts {
		o.apply(u)
	}
	return u
}
func (u *upgraderConfig) Upgrade(w http.ResponseWriter,
	r *http.Request, responseHeader http.Header, opts ...ClientOption) (*Client, error) {

	upgrader := websocket.Upgrader{
		HandshakeTimeout:  u.handshakeTimeout,
		ReadBufferSize:    u.readBufferSize,
		WriteBufferSize:   u.writeBufferSize,
		WriteBufferPool:   u.writeBufferPool,
		Subprotocols:      u.subprotocols,
		Error:             u.errorFunc,
		CheckOrigin:       u.checkOrigin,
		EnableCompression: u.enableCompression,
	}

	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return NewClient(conn, opts...)
}

func (c *Client) ReadErr() error {
	return c.readErr
}

func (c *Client) WriteErr() error {
	return c.writeErr
}

var ErrStopChannelFull = errors.New("stop channel is full")
var ErrNotStarted = errors.New("client has not been started")

func (c *Client) waitStopSend() {
	<-c.stopSend
	close(c.sendDone)
}

func (c *Client) Stop() error {
	select {
	case c.stopSend <- "stop":
	default:
		return ErrStopChannelFull
	}
	return nil
}

func (c *Client) Wait() error {
	c.readWg.Wait()
	return nil
}

func (c *Client) SetMessageHandler(h MessageHandler) {
	c.handler = h
}

func (c *Client) SetID(id string) {
	c.id = id
}

func (c *Client) GetID() string {
	return c.id
}

func (c *Client) SetUserID(id string) {
	c.userID = id
}
func (c *Client) GetUserID() string {
	return c.userID
}

func (c *Client) SetUserGroupID(id string) {
	c.userGroupID = id
}
func (c *Client) GetUserGroupID() string {
	return c.userGroupID
}

func (c *Client) SetHub(h sicore.Hub) {
	c.hub = h
}

func (c *Client) appendReaderOpt(ro sicore.ReaderOption) {
	c.readOpts = append(c.readOpts, ro)
}

var ErrDataChannelClosed = errors.New("send data channel closed")

type msg struct {
	data []byte
	err  chan error
}

func newMsg(data []byte) *msg {
	return &msg{
		data: data,
		err:  make(chan error, 1),
	}
}

func (c *Client) Send(b []byte) error {
	select {
	case <-c.sendDone:
		return ErrDataChannelClosed
	default:
	}

	select {
	case <-c.sendDone:
		return ErrDataChannelClosed
	case c.data <- b:
	}
	return nil
}

func (c *Client) SendAndWait(b []byte) error {
	select {
	case <-c.sendDone:
		return ErrDataChannelClosed
	default:
	}

	m := newMsg(b)
	select {
	case <-c.sendDone:
		return ErrDataChannelClosed
	case c.msg <- m:
		return <-m.err
	}
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

func (c *Client) readPump() {
	defer func() {
		c.Stop()
		c.conn.Close()
		c.hub.Remove(c)
		// if err := c.hub.removeClient(c); err != nil {
		// log.Println("failed to remove client from hub", c.id)
		// }

		// log.Println("return readPump", c.id)
		c.readWg.Done()
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
		if err := c.handler.Handle(r); err != nil {
			c.readErr = err
			return
		}

		cnt++
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(c.pingPeriod)
	normalClose := false
	defer func() {
		ticker.Stop()

		if !normalClose {
			// log.Println("abnormal write")
			c.conn.Close()
		}
		// c.conn.Close() // TODO: should be closed here?
		// Closing here causes losing messages sent by a peer. It is best to close the connection
		// on readPump so it read as much as possible.

		// log.Println("return writePump", c.id)
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
		case message := <-c.msg:
			c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				// log.Println(err)
				c.writeErr = err
				message.err <- err
				return
			}
			w.Write(message.data)
			if err := w.Close(); err != nil {
				// log.Println(err)
				c.writeErr = err
				message.err <- err
				return
			}
			message.err <- nil
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
