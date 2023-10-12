package main

import (
	amqpconsumer "github.com/flylib/mq-consumer"
	"github.com/flylib/mq-consumer/rabbitmq"
	"log"
)

func main() {
	ctx := amqpconsumer.NewContext(
		amqpconsumer.RegisterTopicHandler(new(test)),
	)
	consumer, err := rabbitmq.Dial(ctx, "amqp://admin:admin@192.168.119.128:5672")
	if err != nil {
		panic(err)
	}
	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
