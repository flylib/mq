package mq

const (
	MIME_Type_Binary   = "application/binary"
	MIME_Type_Xml      = "application/xml"
	MIME_Type_Json     = "application/json"
	MIME_Type_Protobuf = "application/x-protobuf"
)

type Option func(ctx *AppContext)

func WithLogger(logger ILogger) Option {
	return func(ctx *AppContext) {
		ctx.ILogger = logger
	}
}

func SetDefaultCodec(codec ICodec) Option {
	return func(ctx *AppContext) {
		ctx.ICodec = codec
	}
}

func RegisterTopicHandler(codec ...ITopicHandler) Option {
	return func(ctx *AppContext) {
		ctx.RegisterTopicHandler(codec...)
	}
}
