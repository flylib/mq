package stream

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type Option func(o *option)

type option struct {
	group, consumer, readMsgIndex string
	maxTryReconnectTimes          uint32
	reconnectionInterval          time.Duration
	heartbeatInterval             time.Duration
	redis.Options
}

func ConsumeGroupInfo(group, name, readMsgIndex string) Option {
	return func(o *option) {
		o.group = group
		o.consumer = name
		o.readMsgIndex = readMsgIndex
	}
}

// less than 1s uses the server's interval
func HeartbeatInterval(duration time.Duration) Option {
	return func(o *option) {
		o.heartbeatInterval = duration
	}
}

// default 10s,reconnection interval
func ReconnectInterval(duration time.Duration) Option {
	return func(o *option) {
		o.reconnectionInterval = duration
	}
}

// default 10,0 means no limit,The maximum number of reconnections allowed after disconnection
func MaxTryReconnectTimes(times uint32) Option {
	return func(o *option) {
		o.maxTryReconnectTimes = times
	}
}

func RedisDB(db int) Option {
	return func(options *option) {
		options.DB = db
	}
}

func RedisAuth(user, pwd string) Option {
	return func(options *option) {
		options.Username = user
		options.Password = pwd
	}
}
