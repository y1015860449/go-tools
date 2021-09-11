package rocketmq

import (
	"context"
	"errors"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/google/uuid"
	"github.com/y1015860449/go-tools/hy_mq/broker"
	"sync"
)

type rocketmqBroker struct {
	addrs []string

	p rocketmq.Producer

	sc []rocketmq.PushConsumer

	connected bool
	scMutex   sync.RWMutex
	opts      broker.Options
}

type subscriber struct {
	t    string
	opts broker.SubscribeOptions
	c    rocketmq.PushConsumer
}

type publication struct {
	c   rocketmq.PushConsumer
	m   *broker.Message
	t   string
	err error
}

func (p *publication) Topic() string {
	return p.t
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (p *publication) Ack() error {
	return nil
}

func (p *publication) Error() error {
	return p.err
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	return s.c.Shutdown()
}

func (r *rocketmqBroker) Address() string {
	if len(r.addrs) > 0 {
		return r.addrs[0]
	}
	return "127.0.0.1:9876"
}

func (r *rocketmqBroker) Connect() error {
	if r.isConnected() {
		return nil
	}

	ropts := make([]producer.Option, 0)

	ropts = append(ropts, producer.WithNsResolver(primitive.NewPassthroughResolver(r.opts.Addrs)))

	if retry, ok := r.opts.Context.Value(retryKey{}).(int); ok {
		ropts = append(ropts, producer.WithRetry(retry))
	}

	if credentials, ok := r.opts.Context.Value(credentialsKey{}).(Credentials); ok {
		ropts = append(ropts, producer.WithCredentials(primitive.Credentials{
			AccessKey: credentials.AccessKey,
			SecretKey: credentials.SecretKey,
		}))
	}

	p, err := rocketmq.NewProducer(ropts...)
	if err != nil {
		return err
	}

	err = p.Start()
	if err != nil {
		return err
	}

	r.scMutex.Lock()
	r.p = p
	r.sc = make([]rocketmq.PushConsumer, 0)
	r.connected = true
	r.scMutex.Unlock()

	return nil
}

func (r *rocketmqBroker) Disconnect() error {
	if !r.isConnected() {
		return nil
	}
	r.scMutex.Lock()
	defer r.scMutex.Unlock()
	for _, consumer := range r.sc {
		consumer.Shutdown()
	}

	r.sc = nil
	r.p.Shutdown()

	r.connected = false
	return nil
}

func (r *rocketmqBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	var cAddrs []string
	for _, addr := range r.opts.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:9876"}
	}
	r.addrs = cAddrs
	return nil
}

func (r *rocketmqBroker) isConnected() bool {
	r.scMutex.RLock()
	defer r.scMutex.RUnlock()
	return r.connected
}

func (r *rocketmqBroker) Options() broker.Options {
	return r.opts
}

func (r *rocketmqBroker) Publish(topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	if !r.isConnected() {
		return errors.New("[rocketmq] broker not connected")
	}

	options := broker.PublishOptions{}
	for _, o := range opts {
		o(&options)
	}

	var (
		delayTimeLevel int
	)

	if options.Context != nil {
		if v, ok := options.Context.Value(delayTimeLevelKey{}).(int); ok {
			delayTimeLevel = v
		}
	}

	m := primitive.NewMessage(topic, msg.Body)

	for k, v := range msg.Header {
		m.WithProperty(k, v)
	}

	if delayTimeLevel > 0 {
		m.WithDelayTimeLevel(delayTimeLevel)
	}
	_, err := r.p.SendSync(context.Background(), m)

	return err
}

func (r *rocketmqBroker) getPushConsumer(groupName string) (rocketmq.PushConsumer, error) {
	ropts := make([]consumer.Option, 0)

	ropts = append(ropts, consumer.WithNsResolver(primitive.NewPassthroughResolver(r.opts.Addrs)))
	ropts = append(ropts, consumer.WithConsumerModel(consumer.Clustering))
	if maxReconsumeTimes, ok := r.opts.Context.Value(maxReconsumeTimesKey{}).(int32); ok {
		ropts = append(ropts, consumer.WithMaxReconsumeTimes(maxReconsumeTimes))
	}

	if credentials, ok := r.opts.Context.Value(credentialsKey{}).(Credentials); ok {
		ropts = append(ropts, consumer.WithCredentials(primitive.Credentials{
			AccessKey: credentials.AccessKey,
			SecretKey: credentials.SecretKey,
		}))
	}
	ropts = append(ropts, consumer.WithGroupName(groupName))

	cs, err := rocketmq.NewPushConsumer(ropts...)
	if err != nil {
		return nil, err
	}

	r.scMutex.Lock()
	defer r.scMutex.Unlock()
	r.sc = append(r.sc, cs)
	return cs, nil
}

func (r *rocketmqBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
		Queue:   uuid.New().String(), // 默认的队列名，会被覆盖的
	}
	for _, o := range opts {
		o(&opt)
	}

	// theoretically, groupName not queue
	// in rocket. one topic have many queue, one queue only belongs to one consumer, one consumer can consume many queue
	// many consumer belongs to a specified group shares messages of the topic
	groupName := opt.Queue
	if len(groupName) == 0 {
		return nil, errors.New("rocketmq need groupName or queue")
	}

	c, err := r.getPushConsumer(groupName)
	if err != nil {
		return nil, err
	}

	err = c.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			rlog.Info("Subscribe process message", map[string]interface{}{
				"Topic":                     msg.Topic,
				"ID":                        msg.MsgId,
				"BornHost":                  msg.BornHost,
				"QueueOffset":               msg.QueueOffset,
				"CommitLogOffset":           msg.CommitLogOffset,
				"PreparedTransactionOffset": msg.PreparedTransactionOffset,
				"ReconsumeTimes":            msg.ReconsumeTimes,
			})

			header := make(map[string]string)
			for k, v := range msg.GetProperties() {
				header[k] = v
			}

			m := &broker.Message{
				Header: header,
				Body:   msg.Body,
			}

			p := &publication{c: c, m: m, t: msg.Topic}
			p.err = handler(p)
			if p.err != nil {
				return consumer.ConsumeRetryLater, p.err
			}
		}

		return consumer.ConsumeSuccess, nil
	})

	err = c.Start()
	if err != nil {
		return nil, err
	}

	return &subscriber{t: topic, opts: opt, c: c}, nil
}

func (r *rocketmqBroker) BrokerName() string {
	return "rocketmq"
}

func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	var cAddrs []string
	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:9876"}
	}

	return &rocketmqBroker{
		addrs: cAddrs,
		opts:  options,
	}
}
