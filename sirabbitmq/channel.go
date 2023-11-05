package sirabbitmq

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultReinitDelay     = 2 * time.Second
	defaultResendDelay     = 5 * time.Second
	defaultConsumeDelay    = 1 * time.Second
	defaultConsumeMaxRetry = 30
)

type Channel struct {
	id string

	conn            *Conn
	channel         *amqp.Channel
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation

	isConnected bool
	done        chan bool
	ready       chan bool

	prefetch     int
	prefetchSize int
	global       bool

	failCount int
}

// NewChannel creates a new consumer state instance, and automatically
// attempts to connect to the server.
func NewChannel(conn *Conn) *Channel {

	c := Channel{
		id:           generateId(),
		conn:         conn,
		done:         make(chan bool),
		ready:        make(chan bool),
		prefetch:     1,
		prefetchSize: 0,
		global:       false,
		failCount:    0,
	}
	go c.handleReinit()
	c.waitReady()
	return &c
}
func NewChannelWithPrefetch(conn *Conn, prefetch int) *Channel {

	c := Channel{
		id:           generateId(),
		conn:         conn,
		done:         make(chan bool),
		ready:        make(chan bool),
		prefetch:     prefetch,
		prefetchSize: 0,
		global:       false,
		failCount:    0,
	}
	go c.handleReinit()
	c.waitReady()
	return &c
}

// handleReconnect will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (c *Channel) handleReinit() bool {
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
func (c *Channel) init() error {
	ch, err := c.conn.GetConnection().Channel()
	if err != nil {
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		return err
	}

	c.changeChannel(ch)

	return nil
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (c *Channel) changeChannel(channel *amqp.Channel) {
	c.channel = channel
	c.notifyChanClose = make(chan *amqp.Error, 1)
	c.channel.NotifyClose(c.notifyChanClose)

	c.notifyConfirm = make(chan amqp.Confirmation, 1)
	c.channel.NotifyPublish(c.notifyConfirm)
}

// Close will cleanly shut down the channel and connection.
func (c *Channel) Close() error {
	if !c.isConnected {
		return errAlreadyClosed
	}
	close(c.done)
	err := c.channel.Close()
	if err != nil {
		return err
	}
	c.isConnected = false
	Infof("closing channel, %s\n", c.id)
	return nil
}

func (c *Channel) GetChannel() *amqp.Channel {
	return c.channel
}

func (c *Channel) waitReady() {
	<-c.ready
}

func (c *Channel) DeclareQueue(queueName string) (amqp.Queue, error) {
	return c.channel.QueueDeclare(queueName, true, false, false, false, nil)
}
func (c *Channel) DeclareOneTimeQueue(queueName string) (amqp.Queue, error) {
	args := make(map[string]interface{})
	args["x-expires"] = 60000
	return c.channel.QueueDeclare(queueName, true, false, false, false, args)
}

func (c *Channel) PushOnce(ctx context.Context, queueName string, data []byte) error {
	_, err := c.DeclareOneTimeQueue(queueName)
	if err != nil {
		return err
	}
	return c.Push(ctx, queueName, data)
}

// Push will push data onto the queue, and wait for a confirm.
// This will block until the server sends a confirm. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (c *Channel) Push(ctx context.Context, queueName string, data []byte) error {
	if !c.isConnected {
		return errNotConnected
	}

	for {
		err := c.push(ctx, queueName, data)
		if err != nil {
			Error("failed to push")
			select {
			case <-c.conn.GetDone():
				return errShutdown
			case <-time.After(defaultResendDelay):
				Error("attempting to push again")
			}
			continue
		}
		confirm := <-c.notifyConfirm
		if confirm.Ack {
			Infof("push has been confirmed(delivery tag: %d)\n", confirm.DeliveryTag)
			return nil
		}
	}
}

func (c *Channel) PushWithReplyTo(ctx context.Context, queueName string, replyTo string, data []byte) error {
	if !c.isConnected {
		return errNotConnected
	}

	for {
		err := c.pushWithReplyTo(ctx, queueName, replyTo, data)
		if err != nil {
			Error("failed to push")
			select {
			case <-c.conn.GetDone():
				return errShutdown
			case <-time.After(defaultResendDelay):
				Error("attempting to push again")
			}
			continue
		}
		confirm := <-c.notifyConfirm
		if confirm.Ack {
			Infof("push has been confirmed(delivery tag: %d)\n", confirm.DeliveryTag)
			return nil
		}
	}
}

// push will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (c *Channel) push(ctx context.Context, queueName string, data []byte) error {
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

func (c *Channel) pushWithReplyTo(ctx context.Context, queueName, replyTo string, data []byte) error {
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

// Consume will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (c *Channel) ConsumeAck(queueName string) (<-chan amqp.Delivery, error) {
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

func (c *Channel) Consume(queueName string) (<-chan amqp.Delivery, error) {
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

func (c *Channel) ConsumeOnce(ctx context.Context, queueName string) ([]byte, error) {
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

type MessageHandler interface {
	Handle(route string, msg []byte) error
}

func (c *Channel) ConsumeWithMessageHandler(ctx context.Context, queueName string, handler MessageHandler) error {
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
				c.failCount++
				if c.failCount > defaultConsumeMaxRetry {
					return errors.New("max retry has been reached.\n" + amqpErr.Error() + "\n" + err.Error())
				}
				Error("failed to consume, trying again... " + err.Error())
				<-time.After(defaultConsumeDelay)
				continue
			}

			// Re-set channel to receive notifications
			// The library closes this channel after abnormal shutdown
			closeChannel = make(chan *amqp.Error, 1)
			c.channel.NotifyClose(closeChannel)

		case delivery := <-deliveries:
			route := delivery.ReplyTo
			if err := handler.Handle(route, delivery.Body); err != nil {
				Errorf("failed to handle message: %s\n", err.Error())
			}
			if err := delivery.Ack(false); err != nil {
				Errorf("failed to acknowledging message: %s\n", err)
			}
		}
	}
}
