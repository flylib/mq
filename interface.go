package mq

type IProducer interface {
	Push(v any) error
}

type IConsumer interface {
	WorkingOn(url string) error
}

type ITopicHandler interface {
	Name() string
	Handler(IMessage)
	OnPanic(IMessage, error)
}

type IMessage interface {
	Ack() error
	//Requeued to be consumed again
	Requeue() error
	//Reject to consume again
	Reject() error
	//Unmarshl To target
	Unmarshal(v any) error
	// MIME content type
	ContentType() string
	// Body raw data
	Body() []byte
	//Message ID
	ID() string
}

// message Codec
type ICodec interface {
	Marshal(v any) (data []byte, err error)
	Unmarshal(data []byte, v any) error
}
