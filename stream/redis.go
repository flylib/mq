package stream

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func connectRedis(host string, o option) (*redis.Client, error) {
	o.Addr = host
	rdb := redis.NewClient(&o.Options)
	return rdb, rdb.Ping(context.Background()).Err()
}
