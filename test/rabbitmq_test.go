package test

import (
	"fmt"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/interface/mq"
	"github.com/flylib/mq/rabbitmq"
	"github.com/flylib/pkg/log/builtinlog"
	"testing"
	"time"
)

type Msg struct {
	Content string `json:"content"`
}

func TestProducer(t *testing.T) {
	broker := rabbitmq.NewBroker(
		rabbitmq.WithDeclareQueues("test"),

		rabbitmq.MustWithLogger(builtinlog.NewLogger()),
		rabbitmq.MustWithCodec(&json.Codec{}),
	)

	err := broker.Connect("amqp://admin:admin@192.168.119.128:5672")
	if err != nil {
		t.Fatal(err)
	}

	////consumer-1 on the same channel
	//err = broker.Subscribe("test", func(message mq.IMessage) {
	//	//time.Sleep(time.Second * 3)
	//	var data Msg
	//	err = message.Unmarshal(&data)
	//	if err != nil {
	//		t.Fatal(err)
	//		return
	//	}
	//	t.Log("[Test1] recvce msg:", data.Content)
	//	message.Ack()
	//	return
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	////consumer-2 on the same channel
	//err = broker.Subscribe("test", func(message mq.IMessage) {
	//	//time.Sleep(time.Second * 3)
	//	var data Msg
	//	err = message.Unmarshal(&data)
	//	if err != nil {
	//		t.Fatal(err)
	//		return
	//	}
	//	t.Log("[Test2] recvce msg:", data.Content)
	//	message.Ack()
	//	return
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	////consumer-3 on other channel
	//
	//channel, err := broker.OpenChannel("")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//err = channel.Subscribe("test", func(message mq.IMessage) {
	//	//time.Sleep(time.Second * 3)
	//	var data Msg
	//	err = message.Unmarshal(&data)
	//	if err != nil {
	//		t.Fatal(err)
	//		return
	//	}
	//	t.Log("[Test3] recvce msg:", data.Content)
	//	message.Ack()
	//	return
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}

	ticker := time.NewTicker(time.Second * 1)
	var i int
	for range ticker.C {
		i++
		msg := Msg{
			Content: fmt.Sprintf("hello-%d", i),
		}
		t.Log("send msg-", i)
		err := broker.Publish("test", msg)
		if err != nil {
			t.Log(err)
		}
	}
}

func TestConsumer(t *testing.T) {
	broker := rabbitmq.NewBroker(
		rabbitmq.WithDeclareQueues("test"),
		rabbitmq.WithPanicHandler(func(message mq.IMessage, err error) {
			t.Fatal(err)
		}),
		rabbitmq.MustWithLogger(builtinlog.NewLogger()),
		rabbitmq.MustWithCodec(&json.Codec{}),
	)
	err := broker.Connect("amqp://admin:admin@192.168.119.128:5672")
	if err != nil {
		t.Fatal(err)
	}
	var i int
	err = broker.Subscribe("test", func(message mq.IMessage) {
		i++
		time.Sleep(time.Second * 3)
		var data Msg
		err = message.Unmarshal(&data)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log("[Test] recvce msg:", data.Content)
		if i == 1 {
			panic("#######")
		}

		message.Ack()
		return
	})
	if err != nil {
		t.Fatal(err)
	}
	//err = broker.Subscribe("test", func(message mq.IMessage) {
	//	i++
	//	time.Sleep(time.Second * 2)
	//	var data Msg
	//	err = message.Unmarshal(&data)
	//	if err != nil {
	//		t.Fatal(err.Error())
	//	}
	//	t.Log("[Test] recvce msg:", data.Content)
	//	message.Ack()
	//	return
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	time.Sleep(time.Hour)
}
