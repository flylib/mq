package rabbitmq

import (
	"github.com/flylib/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"runtime/debug"
	"sync"
	"time"
)

type consumer struct {
	ctx                   *mq.AppContext
	option                option
	conn                  *amqp.Connection
	restartTopicHandlerCh chan mq.ITopicHandler
	reconnecting          sync.Once
	reconnectTimes        uint32
	url                   string
}

func NewConsumer(ctx *mq.AppContext, options ...Option) mq.IConsumer {
	var c = consumer{
		ctx: ctx,
		option: option{
			reconnectionInterval: time.Second * 15,
			maxTryReconnectTimes: 10,
		},
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

	//reset
	c.url = url
	c.reconnectTimes = 0
	c.restartTopicHandlerCh = make(chan mq.ITopicHandler)
	c.reconnecting = sync.Once{}

	//topic channel consume
	c.ctx.RangeTopicHandler(func(topic mq.ITopicHandler) {
		errRun := c.consuming(topic)
		if errRun != nil {
			err = err
			return
		}
	})
	if err != nil {
		return
	}
	for topic := range c.restartTopicHandlerCh {
		c.consuming(topic)
	}

	// reconnect
	for {
		time.Sleep(c.option.reconnectionInterval)
		c.reconnectTimes++
		c.ctx.Infof("Try to reconnect %d times", c.reconnectTimes)
		err = c.Working(url)
		if err != nil {
			c.ctx.Error("reconnect err:", err)
		}
		if c.option.maxTryReconnectTimes != 0 &&
			c.reconnectTimes >= c.option.maxTryReconnectTimes {
			break
		}
	}
	return err
}

func (c *consumer) consuming(handler mq.ITopicHandler) (err error) {
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
				c.ctx.Error("panic error:%s\n\n%s", err, string(debug.Stack()))
				if !c.conn.IsClosed() {
					c.restartTopicHandlerCh <- handler
				}
			} else {
				//Enter reconnection state
				if c.conn.IsClosed() {
					c.reconnecting.Do(func() {
						c.ctx.Error("connection is closed!!!")
						close(c.restartTopicHandlerCh)
					})
				}
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
