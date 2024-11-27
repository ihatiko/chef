package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
)

const (
	component = "kafka"
)

type Client struct {
	cfg    Config
	topics []string
	err    error
}

func (c Client) Name() string {
	return fmt.Sprintf("name: %s topic:%s hosts: %v", component, strings.Join(c.topics, ","), c.cfg.Brokers)
}

func (c Client) Live(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(len(c.cfg.Brokers))
	mt := sync.Mutex{}
	var liveResult []error
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
