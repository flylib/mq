package amqpconsumer

type OnErrorAction int8

type IClient interface {
	Start() error
}

type ITopicHandler interface {
	Topic() string
	OnErrorAction() OnErrorAction
	Handler(IMessage) error
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
}

// message Codec
type ICodec interface {
	MIMEType() string
	Marshal(v any) (data []byte, err error)
	Unmarshal(data []byte, v any) error
}
