package main

import (
	"github.com/flylib/mq"
	"github.com/flylib/mq/rabbitmq"
)

func main() {
	ctx := mq.NewContext(mq.RegisterTopicHandler(new(test)))
	err := rabbitmq.NewConsumer(ctx, rabbitmq.ConsumerName("consumer-test")).WorkingOn("amqp://admin:admin@192.168.119.128:5672")
	if err != nil {
		ctx.Fatal("app exit!!! error:", err)
	}
}
