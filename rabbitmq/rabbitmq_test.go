package rabbitmq

import (
	"fmt"
	"github.com/flylib/mq"
	"log"
	"testing"
	"time"
)

type Msg struct {
	Content string `json:"content"`
}

func TestProducer(t *testing.T) {
	ctx := mq.NewContext()
	broker, err := NewBroker(ctx, "amqp://admin:admin@192.168.119.128:5672")
	if err != nil {
		t.Fatal(err)
	}
	err = broker.DeclareQueue("test")
	if err != nil {
		t.Fatal(err)
	}

	producer, err := broker.NewProducer("")
	if err != nil {
		t.Fatal(err)
	}
	ticker := time.NewTicker(time.Second * 3)
	var i int
	for range ticker.C {
		i++
		msg := Msg{
			Content: fmt.Sprintf("hello-%d", i),
		}
		t.Log("send msg-", i)
		err = producer.Push("test", msg)
		if err != nil {
			t.Fatal()
		}
	}
}

func TestConsumer(t *testing.T) {
	ctx := mq.NewContext(mq.WithRegisterTopicHandler(new(test)))
	err := NewConsumer(ctx,
		WithConsumerName("consumer-test"),
	).Subscribe("amqp://admin:admin@192.168.119.128:5672")
	if err != nil {
		ctx.Fatal("app exit!!! error:", err)
	}
}

type test struct {
}

func (t *test) Topic() string {
	return "test"
}

func (t *test) OnPanic(msg mq.IMessage, err error) {
	fmt.Println(err.Error())
}

func (t *test) Process(msg mq.IMessage) {
	time.Sleep(time.Second * 3)
	var data Msg
	err := msg.Unmarshal(&data)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	log.Println("[Test] recvce msg:", data.Content)
	msg.Ack()
	//panic("panic test")
	return
}
