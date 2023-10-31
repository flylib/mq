package rabbitmq

import (
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
	"github.com/flylib/interface/mq"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type Broker struct {
	amqp.Config
	conn                  *amqp.Connection
	restartTopicHandlerCh chan *Channel
	reconnecting          sync.Once
	reconnectTimes        uint32
	url                   string
	defaultCh             mq.IChannel
	defaultExchange       string
	once                  sync.Once
	reconnectCh           chan bool
	reconnectionInterval  time.Duration
	channels              []*Channel
	consumerName          string
	maxTryReconnectTimes  uint32
	ilog.ILogger          //default logger
	codec.ICodec          //default codec
	sync.Mutex
	declareQueues []string
}

func NewBroker(options ...Option) mq.IBroker {
	b := Broker{
		url:                   "amqp://admin:admin@localhost:5672",
		restartTopicHandlerCh: make(chan *Channel),
		reconnectCh:           make(chan bool),
		reconnectionInterval:  time.Second * 15,
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

func (b *Broker) Connect() (err error) {
	b.conn, err = amqp.Dial(b.url)
	//default channle
	ch, err := b.OpenChannel(b.defaultExchange)
	if err != nil {
		return err
	}
	b.defaultCh = ch
	declarer := b.defaultCh.(interface{ DeclareQueue(queue string) error })
	for _, queue := range b.declareQueues {
		err = declarer.DeclareQueue(queue)
		if err != nil {
			return err
		}
	}

	go func() {
		for item := range b.restartTopicHandlerCh {
			item.Subscribe(item.topic, item.handler)
		}
	}()

	go func() {
		for range b.reconnectCh {
			for {
				err = b.reconnect()
				if err != nil {
					b.ILogger.Error("reconnect err:", err)
					time.Sleep(b.reconnectionInterval)
					continue
				}
				break
			}
		}
	}()

	return
}

func (b *Broker) Close() error {
	return b.conn.Close()
}

func (b *Broker) Publish(topic string, v any) error {
	return b.defaultCh.Publish(topic, v)
}

func (b *Broker) Subscribe(topic string, handler mq.MessageHandler) error {
	return b.defaultCh.Subscribe(topic, handler)
}

func (b *Broker) OpenChannel(name string) (mq.IChannel, error) {
	b.Lock()
	defer b.Unlock()
	channel, err := b.conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Channel{
		Broker:   b,
		ch:       channel,
		exchange: name,
	}, nil
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
		err = b.channels[i].Subscribe(b.channels[i].topic, b.channels[i].handler)
		if err != nil {
			return err
		}
	}
	return nil
}
