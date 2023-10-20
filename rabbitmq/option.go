package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Option func(o *option)

type option struct {
	consumerName         string
	maxTryReconnectTimes uint32
	reconnectionInterval time.Duration
	amqp.Config
}

func WithVhost(vhost string) Option {
	return func(o *option) {
		o.Vhost = vhost
	}
}

func WithConsumerName(name string) Option {
	return func(o *option) {
		o.consumerName = name
	}
}

// less than 1s uses the server's interval
func WithHeartbeatInterval(duration time.Duration) Option {
	return func(o *option) {
		o.Heartbeat = duration
	}
}

// default 10s,reconnection interval
func WithReconnectInterval(duration time.Duration) Option {
	return func(o *option) {
		o.reconnectionInterval = duration
	}
}

// default 10,0 means no limit,The maximum number of reconnections allowed after disconnection
func WithMaxTryReconnectTimes(times uint32) Option {
	return func(o *option) {
		o.maxTryReconnectTimes = times
	}
}
