package main

import (
	"fmt"
	"github.com/flylib/mq"
	"log"
	"time"
)

type test struct {
}

type Msg struct {
	Content string `json:"content"`
}

func (t *test) Name() string {
	return "test"
}

func (t *test) OnPanic(msg mq.IMessage, err error) {
	fmt.Println(err.Error())
}

func (t *test) Handler(msg mq.IMessage) {
	time.Sleep(time.Second * 3)
	var data Msg
	err := msg.Unmarshal(&data)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	log.Println("[Test] recvce msg:", data.Content)
	msg.Ack()
	panic("panic test")
	return
}
