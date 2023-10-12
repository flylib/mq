package amqpconsumer

const (
	NotProcessed OnErrorAction = iota
	Requeue
	Reject
)

const (
	MIMEType_Binary   = "application/binary"
	MIMEType_Xml      = "application/xml"
	MIMEType_Json     = "application/json"
	MIMEType_protobuf = "application/x-protobuf"
)

type Option func(ctx *AppContext)

func RegisterTopicHandler(handler ...ITopicHandler) Option {
	return func(ctx *AppContext) {
		ctx.topicHandlers = append(ctx.topicHandlers, handler...)
	}
}

func RegisterCodec(codec ...ICodec) Option {
	return func(ctx *AppContext) {
		for _, c := range codec {
			ctx.codecs[c.MIMEType()] = c
		}
	}
}
