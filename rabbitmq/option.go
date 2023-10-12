package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Option func(config *option)

type option struct {
	consumerName string
	amqp.Config
}

func UseVhost(vhost string) Option {
	return func(config *option) {
		config.Vhost = vhost
	}
}

// less than 1s uses the server's interval
func HeartbeatInterval(duration time.Duration) Option {
	return func(config *option) {
		config.Heartbeat = duration
	}
}
