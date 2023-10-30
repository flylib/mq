package mq

import (
	"context"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/goutils/logger/log"
	"github.com/flylib/interface/codec"
)

type AppContext struct {
	context.Context
	topicHandlers []IMessageHandler
	ILogger       //default logger
	codec.ICodec  //default codec
}

func NewContext(options ...Option) *AppContext {
	ctx := AppContext{
		Context: context.Background(),
		ILogger: log.NewLogger(log.SyncConsole()),
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
