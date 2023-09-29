package sirabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type UnsafeChannel struct {
	id string

	conn            *Conn
	channel         *amqp.Channel
	notifyChanClose chan *amqp.Error

	isConnected bool
	done        chan bool
	ready       chan bool

	prefetch     int
	prefetchSize int
	global       bool
}

// NewChannel creates a new consumer state instance, and automatically
// attempts to connect to the server.
func NewUnsafeChannel(conn *Conn) *UnsafeChannel {

	c := UnsafeChannel{
		id:           generateId(),
		conn:         conn,
		done:         make(chan bool),
		ready:        make(chan bool),
		prefetch:     1,
		prefetchSize: 0,
		global:       false,
	}
	// runtime.SetFinalizer(&c, func(c *UnsafeChannel) {
	// 	Info("unsafe channel has been finalized")
	// 	c.Close()
	// })
	go c.handleReinit()
	c.waitReady()
	return &c
}
func NewUnsafeChannelWithPrefetch(conn *Conn, prefetch int) *UnsafeChannel {

	c := UnsafeChannel{
		id:           generateId(),
		conn:         conn,
		done:         make(chan bool),
		ready:        make(chan bool),
		prefetch:     prefetch,
		prefetchSize: 0,
		global:       false,
	}
	go c.handleReinit()
	c.waitReady()
	return &c
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (c *UnsafeChannel) handleReinit() bool {
	for {
		c.isConnected = false
		err := c.init()
		if err != nil {
			Warn("failed to initialize a channel")

			select {
			case <-c.done:
				return true
			case <-c.conn.GetDone():
				return true
			case <-time.After(defaultReinitDelay):
				Warn("attempting to re-initialize a channel")
			}
			continue
		}
		Infof("channel(%s) has been initialized\n", c.id)

		c.isConnected = true
		close(c.ready)
		c.ready = make(chan bool)

		select {
		case <-c.done:
			return true
		case <-c.conn.GetDone():
			return true
		case <-c.notifyChanClose:
			Warnf("channel(%s) has been closed, re-initializing...\n", c.id)
		}
	}
}

// init will initialize channel & declare queue
func (c *UnsafeChannel) init() error {
	ch, err := c.conn.GetConnection().Channel()
	if err != nil {
		return err
	}

	// err = ch.Confirm(false)
	// if err != nil {
	// 	return err
	// }

	c.changeChannel(ch)

	return nil
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (c *UnsafeChannel) changeChannel(channel *amqp.Channel) {
	c.channel = channel
	c.notifyChanClose = make(chan *amqp.Error, 1)
	c.channel.NotifyClose(c.notifyChanClose)
}

// Close will cleanly shut down the channel and connection.
func (c *UnsafeChannel) Close() error {
	if !c.isConnected {
		Warn(errAlreadyClosed.Error())
		return errAlreadyClosed
	}
	close(c.done)
	err := c.channel.Close()
	if err != nil {
		Error("failed to close channel: " + err.Error())
		return err
	}
	c.isConnected = false
	Infof("closing unsafe channel, %s\n", c.id)
	return nil
}

func (c *UnsafeChannel) GetChannel() *amqp.Channel {
	return c.channel
}

func (c *UnsafeChannel) waitReady() {
	<-c.ready
}

func (c *UnsafeChannel) DeclareQueue(queueName string) (amqp.Queue, error) {
	return c.channel.QueueDeclare(queueName, true, false, false, false, nil)
}

func (c *UnsafeChannel) DeclareOneTimeQueue(queueName string) (amqp.Queue, error) {
	args := make(map[string]interface{})
	args["x-expires"] = 60000
	return c.channel.QueueDeclare(queueName, true, false, false, false, args)
}

// Push will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (c *UnsafeChannel) Push(ctx context.Context, queueName string, data []byte) error {
	if !c.isConnected {
		return errNotConnected
	}

	return c.channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			// DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

func (c *UnsafeChannel) PushWithReplyTo(ctx context.Context, queueName, replyTo string, data []byte) error {
	if !c.isConnected {
		return errNotConnected
	}

	return c.channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			// DeliveryMode: amqp.Persistent,
			ReplyTo:     replyTo,
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

func (c *UnsafeChannel) PushOnce(ctx context.Context, queueName string, data []byte) error {
	_, err := c.DeclareOneTimeQueue(queueName)
	if err != nil {
		return err
	}
	return c.Push(ctx, queueName, data)
}

// Consume will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (c *UnsafeChannel) ConsumeAck(queueName string) (<-chan amqp.Delivery, error) {
	if !c.isConnected {
		return nil, errNotConnected
	}

	if err := c.channel.Qos(
		c.prefetch,
		c.prefetchSize,
		c.global,
	); err != nil {
		return nil, err
	}

	return c.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *UnsafeChannel) Consume(queueName string) (<-chan amqp.Delivery, error) {
	if !c.isConnected {
		return nil, errNotConnected
	}

	if err := c.channel.Qos(
		c.prefetch,
		c.prefetchSize,
		c.global,
	); err != nil {
		return nil, err
	}

	return c.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (c *UnsafeChannel) ConsumeOnce(ctx context.Context, queueName string) ([]byte, error) {
	_, err := c.DeclareOneTimeQueue(queueName)
	if err != nil {
		return nil, err
	}

	deliveries, err := c.Consume(queueName)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case m := <-deliveries:
		if err := m.Ack(false); err != nil {
			return nil, err
		}
		if _, err := c.channel.QueueDelete(queueName, false, false, false); err != nil {
			return nil, err
		}
		return m.Body, nil
	}
}

func (c *UnsafeChannel) ConsumeWithMessageHandler(ctx context.Context, queueName string, handler MessageHandler) error {
	deliveries, err := c.Consume(queueName)
	if err != nil {
		return err
	}

	// This channel will receive a notification when a channel closed event
	// happens. This must be different from Client.notifyChanClose because the
	// library sends only one notification and Client.notifyChanClose already has
	// a receiver in handleReconnect().
	// Recommended to make it buffered to avoid deadlocks
	closeChannel := make(chan *amqp.Error, 1)
	c.channel.NotifyClose(closeChannel)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case amqpErr := <-closeChannel:
			// This case handles the event of closed channel e.g. abnormal shutdown
			if amqpErr != nil {
				Errorf("channel has been closed due to: %s\n", amqpErr.Error())
			}

			deliveries, err = c.Consume(queueName)
			if err != nil {
				// If the AMQP channel is not ready, it will continue the loop. Next
				// iteration will enter this case because chClosedCh is closed by the
				// library
				Error("failed to consume, trying again...")
				<-time.After(time.Second * 1)
				continue
			}

			// Re-set channel to receive notifications
			// The library closes this channel after abnormal shutdown
			closeChannel = make(chan *amqp.Error, 1)
			c.channel.NotifyClose(closeChannel)

		case delivery := <-deliveries:
			if err := delivery.Ack(false); err != nil {
				Errorf("failed to acknowledging message: %s\n", err)
			}
			route := delivery.ReplyTo
			if err := handler.Handle(route, delivery.Body); err != nil {
				Errorf("failed to handle message: %s\n", err.Error())
			}
		}
	}
}
