package rabbitmq

import (
	"fmt"
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

func (c *consumer) WorkingOn(url string) (err error) {
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
	err = c.ctx.RangeTopicHandler(func(topic mq.ITopicHandler) error {
		err = c.consuming(topic)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
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
		err = c.WorkingOn(url)
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

func (c *consumer) consuming(topic mq.ITopicHandler) (err error) {
	var ch *amqp.Channel
	ch, err = c.conn.Channel()
	if err != nil {
		return
	}

	var deliveryCh <-chan amqp.Delivery
	deliveryCh, err = ch.Consume(
		topic.Name(),          // queue
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
		var msg message = message{ctx: c.ctx}
		// panic handling
		defer func() {
			ch.Close()
			if err := recover(); err != nil {
				c.ctx.Errorf("panic error:%v >>>>>\t\n%s", err, string(debug.Stack()))
				topic.OnPanic(&msg, fmt.Errorf("%v", err))
				if !c.conn.IsClosed() {
					c.restartTopicHandlerCh <- topic
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
			msg.origin = item
			topic.Handler(&msg)
		}
	}()

	return
}
