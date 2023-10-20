package main

import (
	"github.com/flylib/mq"
	"log"
	"time"
)

type test struct {
}

func (t *test) OnPanic(message mq.IMessage, err error) {
	return
}

type Msg struct {
	Content string `json:"content"`
}

func (t *test) Name() string {
	return "test"
}

func (t *test) Handler(msg mq.IMessage) {
	time.Sleep(time.Second * 3)
	var data Msg
	err := msg.Unmarshal(&data)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	log.Println("[test] recvce msg:", data.Content)
	msg.Ack()
	//panic("panic test")
	//msg.Ack()
	return
}
