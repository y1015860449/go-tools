package rocketmq

import (
	"context"

	"github.com/y1015860449/go-tools/hy_mq/broker"
)

type delayTimeLevelKey struct{}
type retryKey struct{}
type maxReconsumeTimesKey struct{}

type credentialsKey struct{}

type Credentials struct {
	AccessKey string
	SecretKey string
}

// WithDelayTimeLevel set message delay time to consume.
// reference delay level definition: 1s 5s 10s 30s 1m 2m 3m 4m 5m 6m 7m 8m 9m 10m 20m 30m 1h 2h
// delay level starts from 1. for example, if we set param level=1, then the delay time is 1s.
func WithDelayTimeLevel(delayLevel int) broker.PublishOption {
	return func(o *broker.PublishOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, delayTimeLevelKey{}, delayLevel)
	}
}

func WithRetry(retry int) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, retryKey{}, retry)
	}
}

func WithMaxReconsumeTimes(maxReconsumeTimes int32) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, maxReconsumeTimesKey{}, maxReconsumeTimes)
	}
}

func WithCredentials(credentials Credentials) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, credentialsKey{}, credentials)
	}
}
