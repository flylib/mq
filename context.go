package mq

import (
	"context"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/interface/codec"
	ilog "github.com/flylib/interface/log"
	"github.com/flylib/pkg/log/builtinlog"
)

type AppContext struct {
	context.Context
	ilog.ILogger  //default logger
	codec.ICodec  //default codec
	topicHandlers []IMessageHandler
}

func NewContext(options ...Option) *AppContext {
	ctx := AppContext{
		Context: context.Background(),
		ILogger: builtinlog.NewLogger(builtinlog.WithSyncConsole()),
		ICodec:  new(json.Codec),
	}
	for _, f := range options {
		f(&ctx)
	}
	return &ctx
}

func (a *AppContext) RegisterTopicHandler(handlers ...IMessageHandler) {
	a.topicHandlers = append(a.topicHandlers, handlers...)
}

func (a *AppContext) RangeTopicHandler(callback func(handler IMessageHandler) error) error {
	for _, handler := range a.topicHandlers {
		err := callback(handler)
		if err != nil {
			return err
		}
	}
	return nil
}
