package sirabbitmq

import (
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultReconnectDelay = 3 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("client is shutting down")
)

type Conn struct {
	id             string
	addr           string
	reconnectDelay time.Duration

	// logger          *log.Logger
	connection      *amqp.Connection
	done            chan bool
	notifyConnClose chan *amqp.Error

	isReady bool
	ready   chan bool
}

// NewConn creates a new consumer state instance, and automatically
// attempts to connect to the server.
func NewConn(addr string) *Conn {
	conn := Conn{
		id:             generateId(),
		addr:           addr,
		reconnectDelay: defaultReconnectDelay,
		done:           make(chan bool),
		ready:          make(chan bool),
	}
	go conn.handleReconnect(addr)
	conn.waitReady()
	return &conn
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (c *Conn) handleReconnect(addr string) {
	for {
		c.isReady = false
		conn, err := c.connect(addr)
		if err != nil {
			Error("failed to connect")

			select {
			case <-c.done:
				return
			case <-time.After(c.reconnectDelay):
				Error("retrying to connect")
			}
			continue
		}

		Infof("connection(%s) has been initialized\n", c.id)

		c.isReady = true
		close(c.ready)
		c.ready = make(chan bool)

		if done := c.handleClose(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (c *Conn) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	c.changeConnection(conn)
	return conn, nil
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (c *Conn) handleClose(conn *amqp.Connection) bool {
	select {
	case <-c.done:
		return true
	case <-c.notifyConnClose:
		Info("connection has been closed. reconnecting...")
		return false
	}
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (c *Conn) changeConnection(connection *amqp.Connection) {
	c.connection = connection
	c.notifyConnClose = make(chan *amqp.Error, 1)
	c.connection.NotifyClose(c.notifyConnClose)
}

// Close will cleanly shut down the channel and connection.
func (c *Conn) Close() error {
	if !c.isReady {
		return errAlreadyClosed
	}
	close(c.done)

	err := c.connection.Close()
	if err != nil {
		return err
	}

	c.isReady = false
	Infof("closing connection, %s\n", c.id)
	return nil
}

func (c *Conn) GetConnection() *amqp.Connection {
	return c.connection
}

func (c *Conn) GetDone() <-chan bool {
	return c.done
}

func (c *Conn) waitReady() {
	<-c.ready
}
