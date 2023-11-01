package rabbitmq

import (
	"context"
	"fmt"
	"github.com/flylib/interface/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"runtime/debug"
	"sync"
)

// channel represents a task
type Channel struct {
	ctx      *Broker
	ch       *amqp.Channel
	exchange string
	sync.Map
	sync.Mutex
	deliveries []DeliveryInfo
}

type DeliveryInfo struct {
	isClosed   bool
	consumerId string
	topic      string
	handler    mq.MessageHandler
	queue      <-chan amqp.Delivery
}

func (c *Channel) Close() error {
	if c.ch.IsClosed() {
		return nil
	}
	return c.ch.Close()
}

func (c *Channel) Publish(topic string, v any) error {
	body, err := c.ctx.ICodec.Marshal(v)
	if err != nil {
		return err
	}
	err = c.ch.PublishWithContext(
		context.Background(),
		c.exchange,
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType:  c.ctx.ICodec.MIMEType(),
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err == amqp.ErrClosed {
		c.ctx.reconnecting.Do(func() {
			c.ctx.reconnectCh <- true
		})
	}
	return err
}

func (c *Channel) Subscribe(topic string, handler mq.MessageHandler) error {
	var (
		deliveryCh <-chan amqp.Delivery
		consumerId = fmt.Sprintf("%s%d", c.ctx.consumerName, c.ctx.serialNumber())
	)
	deliveryCh, err := c.ch.Consume(
		topic,      // queue
		consumerId, // consumer name
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	info := DeliveryInfo{
		consumerId: consumerId,
		topic:      topic,
		handler:    handler,
		queue:      deliveryCh,
	}
	c.Lock()
	c.deliveries = append(c.deliveries, info)
	c.Unlock()

	go c.Delivering(info)
	return nil
}

func (c *Channel) Delivering(info DeliveryInfo) {
	var msg = message{Broker: c.ctx}

	// panic handling
	defer func() {
		//c.ch.Close()
		if err := recover(); err != nil {
			c.ctx.Errorf("panic error:%v >>>>>\t\n%s", err, string(debug.Stack()))
			if !c.ch.IsClosed() {
				c.Lock()
				for i := 0; i < len(c.deliveries); i++ {
					if c.deliveries[i].consumerId == info.consumerId {
						errCancel := c.ch.Cancel(info.consumerId, true)
						if errCancel != nil {
							c.ctx.Errorf("Cancel consumerId-%s error:%s", info.consumerId, errCancel)
						}
						c.deliveries[i].isClosed = true
					}
				}
				c.Unlock()
				return
			}
		}

		//Enter reconnection state
		if c.ctx.conn.IsClosed() {
			c.ctx.reconnecting.Do(func() {
				c.ctx.Error("connection is closed!!!")
				c.ctx.reconnectCh <- true
			})
		}
	}()

	for item := range info.queue {
		msg.origin = item
		info.handler(&msg)
	}
}

func (c *Channel) DeclareQueue(queue string) error {
	// declare a queue to hold message
	_, err := c.ch.QueueDeclare(
		queue, // 队列名
		true,  // 是否持久化
		false, // 是否自动删除
		false, // 是否排他
		false, // 是否阻塞
		nil,   // 其他参数
	)
	return err
}
