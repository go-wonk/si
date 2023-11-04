package sirabbitmq

import (
	"context"
	"sync"
)

type Consumer struct {
	addr        string
	conn        *Conn
	numChannels int

	channels []*Channel
}

func NewConsumer(addr string, numChannels int, prefetch int) *Consumer {
	consumer := &Consumer{
		addr:        addr,
		numChannels: numChannels,
		channels:    make([]*Channel, numChannels),
	}

	consumer.conn = NewConn(addr, prefetch)
	for i := 0; i < numChannels; i++ {
		consumer.channels[i] = NewChannelWithPrefetch(consumer.conn, prefetch)
	}
	return consumer
}

func (c *Consumer) Close() []error {
	errs := make([]error, 0)
	l := len(c.channels)
	for i := 0; i < l; i++ {
		err := c.channels[i].Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	err := c.conn.Close()
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

func (c *Consumer) ConsumeWithMessageHandler(ctx context.Context, queueName string, handler MessageHandler) error {

	l := len(c.channels)
	wg := sync.WaitGroup{}
	wg.Add(l)
	for i := 0; i < l; i++ {
		go func(wg *sync.WaitGroup, index int, ctx context.Context, handler MessageHandler) {
			defer wg.Done()
			c.channels[index].ConsumeWithMessageHandler(ctx, queueName, handler)
		}(&wg, i, ctx, handler)
	}

	wg.Wait()
	return nil
}

// type DynamicConsumer struct {
// 	ch *Channel
// }

// func NewDynamicConsumer(ch *Channel) *DynamicConsumer {
// 	return &DynamicConsumer{
// 		ch: ch,
// 	}
// }

// func (c *DynamicConsumer) Close() error {
// 	return c.ch.Close()
// }

// func (c *DynamicConsumer) ConsumeWithMessageHandler(ctx context.Context, queueName string, handler MessageHandler) error {
// 	return c.ch.ConsumeWithMessageHandler(ctx, queueName, handler)
// }
