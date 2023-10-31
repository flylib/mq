package main

import (
	"github.com/flylib/interface/mq"
	"github.com/flylib/mq/stream"
)

func main() {
	ctx := mq.NewContext(
		mq.WithRegisterTopicHandler(new(test)),
	)
	err := stream.NewConsumer(ctx).Subscribe("192.168.119.128:6379")
	if err != nil {
		ctx.Fatal("app exit!!! error:", err)
	}
}
