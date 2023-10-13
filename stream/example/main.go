package main

import (
	"github.com/flylib/mq"
	"github.com/flylib/mq/stream"
)

func main() {
	ctx := mq.NewContext(
		mq.RegisterTopicHandler(new(test)),
	)
	err := stream.NewConsumer(ctx).Working("192.168.119.128:6379")
	if err != nil {
		ctx.Fatal("app exit!!! error:", err)
	}
}
