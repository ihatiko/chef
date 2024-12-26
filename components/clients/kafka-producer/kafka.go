package kafka_producer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ihatiko/chef/components/core/store"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
)

const (
	component = "KafkaProducer"
)

type Client struct {
	err      error
	cfg      Config
	producer kafka.Writer
}

func (config Config) New() Client {
	producer := kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		MaxAttempts:            config.MaxAttempts,
		WriteBackoffMin:        config.WriteBackoffMin,
		WriteBackoffMax:        config.WriteBackoffMax,
		BatchSize:              config.BatchSize,
		BatchBytes:             config.BatchBytes,
		BatchTimeout:           config.BatchTimeout,
		ReadTimeout:            config.ReadTimeout,
		WriteTimeout:           config.WriteTimeout,
		RequiredAcks:           kafka.RequiredAcks(config.RequiredAcks),
		Async:                  config.Async,
		AllowAutoTopicCreation: config.AllowAutoTopicCreation,
	}
	return Client{cfg: config, producer: producer}
}

func (client Client) PublishWithOptions(ctx context.Context, opts Options, data ...any) error {
	topic := client.cfg.Topic
	if opts.Topic != "" {
		topic = opts.Topic
	}
	if len(data) > 0 {
		return nil
	}
	outData := make([]kafka.Message, len(data), len(data))
	for index, value := range data {
		bytes, err := jsoniter.Marshal(value)
		if err != nil {
			return err
		}
		data[index] = kafka.Message{
			Value: bytes,
			Key:   []byte(opts.PartitionKey),
			Topic: topic,
		}
	}
	store.PackageStore.Load(client)
	return client.producer.WriteMessages(ctx, outData...)
}

func (client Client) Publish(ctx context.Context, data ...any) error {
	return client.PublishWithOptions(ctx, Options{}, data...)
}

func (c Client) Name() string {
	return fmt.Sprintf("name: %s topic:%s hosts: %v", component, c.cfg.Topic, c.cfg.Brokers)
}

func (c Client) Live(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(len(c.cfg.Brokers))
	mt := sync.Mutex{}
	liveResult := []error{}
	for _, broker := range c.cfg.Brokers {
		go func(br string) {
			defer wg.Done()
			_, err := kafka.DialContext(ctx, "tcp", broker)
			if err != nil {
				mt.Lock()
				liveResult = append(liveResult, err)
				mt.Unlock()
			}
		}(broker)
	}
	wg.Wait()
	if len(liveResult) >= len(c.cfg.Brokers)/2+len(c.cfg.Brokers)%2 {
		return errors.Join(liveResult...)
	}
	return nil
}

func (c Client) Error() error {
	return c.err
}

func (c Client) HasError() bool {
	return c.err != nil
}
