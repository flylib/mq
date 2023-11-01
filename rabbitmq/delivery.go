package rabbitmq

import (
	"github.com/flylib/interface/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Delivery struct {
	isClosed   bool
	consumerId string
	topic      string
	handler    mq.MessageHandler
	queue      <-chan amqp.Delivery
}
