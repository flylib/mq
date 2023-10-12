package amqpconsumer

import (
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/goutils/codec/protobuf"
	"github.com/flylib/goutils/logger/log"
)

type AppContext struct {
	topicHandlers []ITopicHandler
	//message codec,default support json and protobuf codec
	codecs map[string]ICodec
	ILogger
}

func NewContext(options ...Option) *AppContext {
	ctx := AppContext{
		codecs:  make(map[string]ICodec),
		ILogger: log.NewLogger(log.SyncConsole()),
	}
	for _, f := range options {
		f(&ctx)
	}
	return &ctx
}

func (a *AppContext) RangeTopicHandler(callback func(handler ITopicHandler)) {
	for _, handler := range a.topicHandlers {
		callback(handler)
	}
}

var (
	jsonCodec     = new(json.Codec)
	protobufCodec = new(protobuf.Codec)
)

func (a *AppContext) GetCodecByMIMEType(mimeType string) (ICodec, bool) {
	switch mimeType {
	case MIMEType_Json:
		return jsonCodec, true
	case MIMEType_protobuf:
		return protobufCodec, true
	default:
		codec, ok := a.codecs[mimeType]
		return codec, ok
	}
}
