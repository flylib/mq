package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type message struct {
	*Broker
	origin amqp.Delivery
}

func (m *message) ID() string {
	return m.origin.MessageId
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
	return m.ICodec.Unmarshal(m.origin.Body, v)
}

func (m *message) ContentType() string {
	return m.origin.ContentType
}

func (m *message) Body() []byte {
	return m.origin.Body
}
