package rabbitmq

import (
	"github.com/flylib/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"runtime/debug"
)

type consumer struct {
	ctx       *mq.AppContext
	option    option
	conn      *amqp.Connection
	restartCh chan mq.ITopicHandler
}

func NewConsumer(ctx *mq.AppContext, options ...Option) mq.IConsumer {
	var c = consumer{
		ctx:       ctx,
		restartCh: make(chan mq.ITopicHandler),
	}
	for _, f := range options {
		f(&c.option)
	}
	return &c
}

func (c *consumer) Working(url string) (err error) {
	c.conn, err = amqp.DialConfig(url, c.option.Config)
	if err != nil {
		return
	}
	c.ctx.RangeTopicHandler(func(handler mq.ITopicHandler) {
		err = c.working(handler)
		if err != nil {
			err = err
			return
		}
	})

	if err != nil {
		return
	}
	for handler := range c.restartCh {
		c.working(handler)
	}
	return nil
}

func (c *consumer) working(handler mq.ITopicHandler) (err error) {

	var ch *amqp.Channel
	ch, err = c.conn.Channel()
	if err != nil {
		return
	}

	var deliveryCh <-chan amqp.Delivery
	deliveryCh, err = ch.Consume(
		handler.Topic(),       // queue
		c.option.consumerName, // consumer name
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}
	go func() {
		// panic handling
		defer func() {
			ch.Close()
			if err := recover(); err != nil {
				c.ctx.Error("panic error:%s\n%s", err, string(debug.Stack()))
				c.restartCh <- handler
			}
		}()

		for item := range deliveryCh {
			msg := message{
				origin: item,
				ctx:    c.ctx,
			}
			err := handler.Handler(&msg)
			if err != nil {
				switch handler.OnErrorAction() {
				case mq.Reject:
					msg.Reject()
				case mq.Requeue:
					msg.Requeue()
				}
			}
		}
	}()

	return
}
