package stream

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
)

type message struct {
	c      *consumer
	origin redis.XMessage
	stream string
}

func (m *message) ID() string {
	return m.origin.ID
}

func (m *message) Ack() error {
	return m.c.rdb.XAck(m.c.ctx, m.stream, m.c.option.group, m.ID()).Err()
}

func (m *message) Requeue() error {
	return nil
}

func (m *message) Reject() error {
	return m.c.rdb.XAck(m.c.ctx, m.stream, m.c.option.group, m.ID()).Err()
}

func (m *message) Unmarshal(v any) error {
	return mapstructure.Decode(m.origin.Values, v)
}

func (m *message) ContentType() string {
	v, ok := m.origin.Values["content_type"]
	if !ok {
		return ""
	}
	t, _ := v.(string)
	return t
}

func (m *message) Body() []byte {
	marshal, _ := json.Marshal(m.origin.Values)
	return marshal
}
