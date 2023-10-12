package rabbitmq

import (
	"errors"
	"fmt"
	amqpconsumer "github.com/flylib/mq-consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

type message struct {
	origin amqp.Delivery
	ctx    *amqpconsumer.AppContext
}

func (m *message) Ack() error {
	return m.origin.Ack(false)
}

func (m *message) Requeue() error {
	return m.origin.Reject(true)
}

func (m *message) Reject() error {
	return m.origin.Reject(false)
}

func (m *message) Unmarshal(v any) error {
	codec, ok := m.ctx.GetCodecByMIMEType(m.origin.ContentType)
	if !ok {
		return errors.New(fmt.Sprintf("Unsupported parsing '%s' type", m.origin.ContentType))
	}
	return codec.Unmarshal(m.origin.Body, v)
}

func (m *message) ContentType() string {
	return m.origin.ContentType
}

func (m *message) Body() []byte {
	return m.origin.Body
}
