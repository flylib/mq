package rabbitmq

import (
	amqpconsumer "github.com/flylib/mq-consumer"
	amqp "github.com/rabbitmq/amqp091-go"
	"runtime/debug"
)

type consumer struct {
	ctx       *amqpconsumer.AppContext
	option    option
	conn      *amqp.Connection
	restartCh chan amqpconsumer.ITopicHandler
}

func Dial(ctx *amqpconsumer.AppContext, url string, options ...Option) (amqpconsumer.IClient, error) {
	var c = consumer{
		ctx:       ctx,
		restartCh: make(chan amqpconsumer.ITopicHandler),
	}
	for _, f := range options {
		f(&c.option)
	}
	conn, err := amqp.DialConfig(url, c.option.Config)
	c.conn = conn
	return &c, err
}

func (c *consumer) Start() (err error) {
	c.ctx.RangeTopicHandler(func(handler amqpconsumer.ITopicHandler) {
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

func (c *consumer) working(handler amqpconsumer.ITopicHandler) (err error) {

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
				case amqpconsumer.Reject:
					msg.Reject()
				case amqpconsumer.Requeue:
					msg.Requeue()
				}
			}
		}
	}()

	return
}
