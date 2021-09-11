package kafka

import "github.com/micro/go-micro/v2/broker"

const shardKey = "shardKey"

func NewShardMessage(value string, header map[string]string, body []byte) *broker.Message {
	if header == nil && value != "" {
		header = make(map[string]string)
	}
	if header != nil && value != "" {
		header[shardKey] = value
	}
	return &broker.Message{
		Header: header,
		Body:   body,
	}
}
