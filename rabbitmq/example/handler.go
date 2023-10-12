package main

import (
	"github.com/flylib/mq"
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

func (t *test) OnErrorAction() mq.OnErrorAction {
	return mq.NotProcessed
}

func (t *test) Handler(msg mq.IMessage) error {
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
