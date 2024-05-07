package kafka_producer

import "context"

type Options struct {
	Topic        string
	PartitionKey string
}

type IProducer interface {
	PublishWithOptions(ctx context.Context, opts Options, data ...any) error
	Publish(ctx context.Context, data ...any) error
}
