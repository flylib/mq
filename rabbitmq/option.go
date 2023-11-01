package rabbitmq

import (
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
	"time"
)

type Option func(o *Broker)

func WithVhost(vhost string) Option {
	return func(o *Broker) {
		o.Vhost = vhost
	}
}

func WithConsumerName(name string) Option {
	return func(o *Broker) {
		o.consumerName = name
	}
}

// less than 1s uses the server's interval
func WithHeartbeatInterval(duration time.Duration) Option {
	return func(o *Broker) {
		o.Heartbeat = duration
	}
}

// default 10s,reconnection interval
func WithReconnectInterval(duration time.Duration) Option {
	return func(o *Broker) {
		o.reconnectionInterval = duration
	}
}

// default 10,0 means no limit,The maximum number of reconnections allowed after disconnection
func WithMaxTryReconnectTimes(times uint32) Option {
	return func(o *Broker) {
		o.maxTryReconnectTimes = times
	}
}

func WithDefaultExchange(exchange string) Option {
	return func(o *Broker) {
		o.defaultExchange = exchange
	}
}

func WithDeclareQueues(queues ...string) Option {
	return func(o *Broker) {
		o.declareQueues = append(o.declareQueues, queues...)
	}
}

func MustWithLogger(logger ilog.ILogger) Option {
	return func(o *Broker) {
		o.ILogger = logger
	}
}

func MustWithCodec(c codec.ICodec) Option {
	return func(o *Broker) {
		o.ICodec = c
	}
}
