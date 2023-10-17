package stream

import (
	"github.com/flylib/mq"
	"github.com/redis/go-redis/v9"
	"time"
)

type consumer struct {
	ctx                   *mq.AppContext
	rdb                   *redis.Client
	option                option
	readArg               redis.XReadGroupArgs
	restartTopicHandlerCh chan mq.ITopicHandler
}

func NewConsumer(ctx *mq.AppContext, options ...Option) mq.IConsumer {
	var c = consumer{
		ctx:                   ctx,
		restartTopicHandlerCh: make(chan mq.ITopicHandler),
		option: option{
			reconnectionInterval: time.Second * 15,
			maxTryReconnectTimes: 10,
			group:                "group-0",
			consumer:             "consumer-0",
			readMsgIndex:         "$",
		},
	}
	for _, f := range options {
		f(&c.option)
	}
	return &c
}

func (c *consumer) WorkingOn(url string) (err error) {
	c.rdb, err = connectRedis(url, c.option)
	if err != nil {
		return err
	}

	err = c.ctx.RangeTopicHandler(func(stream mq.ITopicHandler) error {
		//err = c.createConsumeGroup(stream)
		//if err != nil {
		//	return err
		//}
		//arg := redis.XReadArgs{
		//	Streams: []string{},
		//	Count:   1,
		//	Block:   0,
		//}
		return c.consuming(stream)
	})
	if err != nil {
		return err
	}

	for stream := range c.restartTopicHandlerCh {
		c.consuming(stream)
	}

	return err
}

func (c *consumer) createConsumeGroup(stream mq.ITopicHandler) error {
	//get groups info
	groups, err := c.rdb.XInfoGroups(c.ctx, stream.Name()).Result()
	if err != nil {
		return err
	}
	var isHaveGroup bool
	for _, item := range groups {
		if item.Name == c.option.group {
			isHaveGroup = true
			break
		}
	}
	if !isHaveGroup {
		_, err = c.rdb.XGroupCreate(c.ctx, stream.Name(), c.option.group, c.option.readMsgIndex).Result()
		if err != nil {
			return err
		}
	}

	//get consumers info
	var isHaveConsumer bool
	consumers, err := c.rdb.XInfoConsumers(c.ctx, stream.Name(), c.option.group).Result()
	if err != nil {
		return err
	}
	for _, item := range consumers {
		if item.Name == c.option.consumer {
			isHaveConsumer = true
			break
		}
	}
	if !isHaveConsumer {
		_, err = c.rdb.XGroupCreateConsumer(c.ctx, stream.Name(), c.option.group, c.option.consumer).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *consumer) consuming(stream mq.ITopicHandler) (err error) {
	err = c.createConsumeGroup(stream)
	if err != nil {
		return err
	}
	arg := redis.XReadGroupArgs{
		Group:    c.option.group,
		Consumer: c.option.consumer,
		Streams:  []string{stream.Name(), c.option.readMsgIndex},
		Count:    1,
		Block:    0,
	}
	go func() {
		var msg message = message{c: c, stream: stream.Name()}
		for {
			var streams []redis.XStream
			streams, err = c.rdb.XReadGroup(c.ctx, &arg).Result()
			if err != nil {
				c.ctx.Errorf("xread [%s] error:%s", stream.Name(), err.Error())
				return
			}
			if len(streams) > 0 {
				for _, item := range streams[0].Messages {
					msg.origin = item
					err = stream.Handler(&msg)
					if err != nil {
						switch stream.OnErrorAction() {
						case mq.Reject:
							msg.Reject()
						case mq.Requeue:
							msg.Requeue()
						}
					}
				}
			}
		}
	}()
	return nil
}
