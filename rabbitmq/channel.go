package rabbitmq

import (
	"context"
	"github.com/flylib/interface/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"runtime/debug"
)

// 一个通道代表着一个任务
type Channel struct {
	*Broker
	ch       *amqp.Channel
	exchange string
	topic    string
	handler  mq.MessageHandler
}

func (c *Channel) Close() error {
	if c.ch.IsClosed() {
		return nil
	}
	return c.ch.Close()
}

func (c *Channel) Publish(topic string, v any) error {
	body, err := c.ICodec.Marshal(v)
	if err != nil {
		return err
	}
	return c.ch.PublishWithContext(
		context.Background(),
		c.exchange,
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType:  c.ICodec.MIMEType(),
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

func (c *Channel) Subscribe(topic string, handler mq.MessageHandler) error {
	var deliveryCh <-chan amqp.Delivery
	deliveryCh, err := c.ch.Consume(
		topic,          // queue
		c.consumerName, // consumer name
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		var msg = message{Broker: c.Broker}

		// panic handling
		defer func() {
			c.ch.Close()
			if err := recover(); err != nil {
				c.ILogger.Errorf("panic error:%v >>>>>\t\n%s", err, string(debug.Stack()))
				if !c.conn.IsClosed() {
					c.restartTopicHandlerCh <- c
				}
			} else {
				//Enter reconnection state
				if c.conn.IsClosed() {
					c.reconnecting.Do(func() {
						c.ILogger.Error("connection is closed!!!")
						close(c.restartTopicHandlerCh)
					})
				}
			}
		}()

		for item := range deliveryCh {
			msg.origin = item
			handler(&msg)
		}
	}()
	return nil
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
