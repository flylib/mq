package main

import (
	"github.com/flylib/mq"
	"github.com/flylib/mq/rabbitmq"
	"log"
)

func main() {
	ctx := mq.NewContext(
		mq.RegisterTopicHandler(new(test)),
	)
	err := rabbitmq.NewConsumer(ctx).Working("amqp://admin:admin@192.168.119.128:5672")
	if err != nil {
		log.Fatal(err)
	}
}
