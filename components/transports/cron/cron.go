package cron

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const (
	defaultTimeout = 10 * time.Second
	defaultWorker  = 1
	componentName  = "cron"
)

type Request struct {
	ctx context.Context
	id  int
}

func (c *Request) Context() context.Context {
	return c.Context()
}

type h func(context Request) error

type transport struct {
	Config *Config
	h      h
}

func (t transport) Name() string {
	return fmt.Sprintf("%s id: %s", componentName, uuid.New().String())
}

type Options func(*transport)

func (cfg *Config) Use(opts ...Options) *transport {
	t := new(transport)

	t.Config = cfg
	for _, opt := range opts {
		opt(t)
	}
	if t.Config.Timeout == 0 {
		t.Config.Timeout = defaultTimeout
	}
	if t.Config.Workers != 0 {
		t.Config.Workers = defaultWorker
	}

	return t
}

func (t transport) Routing(h h) transport {
	t.h = h
	return t
}
func (t transport) Live(ctx context.Context) error {
	return nil
}
func (t transport) Run() {
	if t.h == nil {
		otelzap.L().Fatal("cron transport handler is nil")
	}
	wg := &sync.WaitGroup{}
	for i := range t.Config.Workers {
		wg.Add(1)
		defer wg.Done()
		slog.Info("Start cron worker",
			slog.Int("worker", i+1),
		)
		t.handler(i + 1)
		slog.Info("End cron worker",
			slog.Int("worker", i+1),
		)
	}
	wg.Wait()
}

func (t transport) handler(id int) {
	ctx, cancel := context.WithTimeout(context.TODO(), t.Config.Timeout)
	defer func() {
		if r := recover(); r != nil {
			slog.Error("error handling message", zap.Any("panic", r))
		}
	}()

	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				slog.Warn("context deadline exceeded daemon")
				return
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return
			}
			if ctx.Err() != nil {
				otelzap.S().Error(ctx.Err())
			}
		}
	}()
	err := t.h(Request{ctx: ctx, id: id})
	if err != nil {
		slog.Error("Error daemon worker", slog.String("desc", err.Error()))
	}
}

func (t transport) TimeToWait() time.Duration {
	return t.Config.Timeout
}

func (t transport) Shutdown() error {
	return nil
}
