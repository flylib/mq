package rabbitmq

import (
	"context"
	"github.com/flylib/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	ctx  *mq.AppContext
	conn *amqp.Connection
	url  string
	//option      option
}

func NewBroker(ctx *mq.AppContext, url string) (*Broker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	b := &Broker{
		ctx:  ctx,
		conn: conn,
		url:  url,
	}
	return b, nil
}

type Producer struct {
	exchange string //exhange
	ch       *amqp.Channel
	*Broker
}

func (b *Broker) NewProducer(exchange string) (mq.IProducer, error) {
	channel, err := b.conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Producer{
		exchange: exchange,
		ch:       channel,
		Broker:   b,
	}, nil
}

func (p *Producer) Push(route string, v any) error {
	data, err := p.ctx.Marshal(v)
	if err != nil {
		return err
	}
	return p.ch.PublishWithContext(
		context.Background(),
		p.exchange,
		route,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  p.ctx.MIMEType(),
			Body:         data,
		},
	)
}
