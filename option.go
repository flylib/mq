package mq

import "github.com/flylib/interface/codec"

type Option func(ctx *AppContext)

func WithLogger(logger ILogger) Option {
	return func(ctx *AppContext) {
		ctx.ILogger = logger
	}
}

func WithCodec(codec codec.ICodec) Option {
	return func(ctx *AppContext) {
		ctx.ICodec = codec
	}
}

func WithRegisterTopicHandler(codec ...IMessageHandler) Option {
	return func(ctx *AppContext) {
		ctx.RegisterTopicHandler(codec...)
	}
}
