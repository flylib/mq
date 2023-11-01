package rabbitmq

import (
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
	"github.com/flylib/interface/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"sync/atomic"
	"time"
)

type Broker struct {
	amqp.Config
	conn           *amqp.Connection
	reconnecting   sync.Once
	reconnectTimes uint32
	url            string
	//defaultCh            mq.IChannel
	defaultExchange      string
	once                 sync.Once
	reconnectCh          chan bool
	reconnectionInterval time.Duration
	channels             []*Channel
	consumerName         string
	maxTryReconnectTimes uint32
	ilog.ILogger         //default logger
	codec.ICodec         //default codec
	sync.Mutex
	declareQueues []string
	serial        int32
}

func NewBroker(options ...Option) mq.IBroker {
	b := Broker{
		url:                  "amqp://admin:admin@localhost:5672",
		reconnectCh:          make(chan bool),
		reconnectionInterval: time.Second * 15,
	}
	for _, f := range options {
		f(&b)
	}
	if b.ILogger == nil {
		panic("Ilogger is nil")
	}
	if b.ICodec == nil {
		panic("ICodec is nil")
	}
	return &b

}

func (b *Broker) Connect(url string) (err error) {
	b.conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}
	b.url = url
	//default channel
	defaultCh, err := b.OpenChannel(b.defaultExchange)
	if err != nil {
		return err
	}

	//declare queues for store message
	declarer := defaultCh.(interface{ DeclareQueue(queue string) error })
	for _, queue := range b.declareQueues {
		err = declarer.DeclareQueue(queue)
		if err != nil {
			return err
		}
	}

	go func() {
		for range b.reconnectCh {
			for {
				b.ILogger.Info("try reconnect to ", b.url)
				err = b.reconnect()
				if err != nil {
					b.ILogger.Error("reconnect err:", err)
					time.Sleep(b.reconnectionInterval)
					continue
				}
				b.ILogger.Info("reconnect success!!! ")
				break
			}
			b.reconnecting = sync.Once{}
		}
	}()

	return
}

func (b *Broker) Close() error {
	close(b.reconnectCh)
	return b.conn.Close()
}

func (b *Broker) Publish(topic string, v any) error {
	return b.channels[0].Publish(topic, v)
}

func (b *Broker) Subscribe(topic string, handler mq.MessageHandler) error {
	return b.channels[0].Subscribe(topic, handler)
}

func (b *Broker) OpenChannel(name string) (mq.IChannel, error) {
	b.Lock()
	defer b.Unlock()

	channel, err := b.conn.Channel()
	if err != nil {
		return nil, err
	}
	c := &Channel{
		ctx:      b,
		ch:       channel,
		exchange: name,
	}
	b.channels = append(b.channels, c)
	return c, nil
}

func (b *Broker) reconnect() error {
	b.Lock()
	defer b.Unlock()

	conn, err := amqp.Dial(b.url)
	if err != nil {
		return err
	}
	b.conn = conn

	for i := 0; i < len(b.channels); i++ {
		var channel *amqp.Channel
		channel, err = conn.Channel()
		if err != nil {
			return err
		}
		b.channels[i].ch = channel

		//var deliveries []<-chan amqp.Delivery
		for j, item := range b.channels[i].deliveries {
			if item.isClosed {
				continue
			}
			queue, err := channel.Consume(
				item.topic,      // queue
				item.consumerId, // consumer name
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				return err
			}
			b.channels[i].deliveries[j].queue = queue
			go b.channels[i].Delivering(b.channels[i].deliveries[j])
		}
	}
	return nil
}

func (b *Broker) serialNumber() int32 {
	atomic.AddInt32(&b.serial, 1)
	return b.serial
}
