package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Option func(config *option)

type option struct {
	consumerName string

	maxTryReconnectTimes uint32

	reconnectionInterval time.Duration
	amqp.Config
}

func UseVhost(vhost string) Option {
	return func(config *option) {
		config.Vhost = vhost
	}
}

func ConsumerName(name string) Option {
	return func(config *option) {
		config.consumerName = name
	}
}

// less than 1s uses the server's interval
func HeartbeatInterval(duration time.Duration) Option {
	return func(config *option) {
		config.Heartbeat = duration
	}
}

// default 10s,reconnection interval
func ReconnectInterval(duration time.Duration) Option {
	return func(config *option) {
		config.reconnectionInterval = duration
	}
}

// default 10,0 means no limit,The maximum number of reconnections allowed after disconnection
func MaxTryReconnectTimes(times uint32) Option {
	return func(config *option) {
		config.maxTryReconnectTimes = times
	}
}
