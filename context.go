package mq

import (
	"context"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/goutils/logger/log"
)

type AppContext struct {
	context.Context
	topicHandlers []ITopicHandler
	ILogger       //default logger
	ICodec        //default codec
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

func (a *AppContext) RegisterTopicHandler(handlers ...ITopicHandler) {
	a.topicHandlers = append(a.topicHandlers, handlers...)
}

func (a *AppContext) RangeTopicHandler(callback func(handler ITopicHandler) error) error {
	for _, handler := range a.topicHandlers {
		err := callback(handler)
		if err != nil {
			return err
		}
	}
	return nil
}
