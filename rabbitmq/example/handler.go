package main

import (
	amqpconsumer "github.com/flylib/mq-consumer"
	"log"
	"time"
)

type test struct {
}
type Msg struct {
	Content string `json:"content"`
}

func (t *test) Topic() string {
	return "test"
}

func (t *test) OnErrorAction() amqpconsumer.OnErrorAction {
	return amqpconsumer.NotProcessed
}

func (t *test) Handler(msg amqpconsumer.IMessage) error {
	time.Sleep(time.Second)
	var data Msg
	err := msg.Unmarshal(&data)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	log.Println("[test] recvce msg:", data.Content)
	msg.Ack()
	panic("panic test")
	//msg.Ack()
	return nil
}
