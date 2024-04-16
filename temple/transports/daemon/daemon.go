package daemon

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const (
	defaultInterval = 10 * time.Second
	defaultWorker   = 1
)

type Config struct {
	Timeout  time.Duration
	Interval time.Duration
	Workers  int
}

type Context struct {
	ctx context.Context
}

func (c *Context) Context() context.Context {
	return c.Context()
}

type H func(context Context) error

type Transport struct {
	Config *Config
	h      H
	Ticker *time.Ticker
}
type Options func(*Transport)

func (cfg *Config) Use(opts ...Options) *Transport {
	t := new(Transport)

	t.Config = cfg
	for _, opt := range opts {
		opt(t)
	}
	if t.Config.Interval == 0 {
		t.Ticker = time.NewTicker(defaultInterval)
	} else {
		t.Ticker = time.NewTicker(cfg.Interval)
	}
	if t.Config.Workers != 0 {
		t.Config.Workers = defaultWorker
	}

	return t
}

func Router(router func(h H)) Options {
	return func(t *Transport) {
		router(t.h)
	}
}

func (t *Transport) Run() {
	for range t.Ticker.C {
		wg := &sync.WaitGroup{}
		for i := range t.Config.Workers {
			wg.Add(1)
			defer wg.Done()
			otelzap.S().Infof("Start daemon worker",
				zap.Int("number", i+1),
			)
			t.handler()
		}
		wg.Wait()
	}
}

func (t *Transport) handler() {
	ctx, cancel := context.WithTimeout(context.TODO(), t.Config.Timeout)
	defer func() {
		if r := recover(); r != nil {
			otelzap.
				Ctx(ctx).Error("error handling message", zap.Any("panic", r))
		}
	}()

	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				otelzap.S().Warn("context deadline exceeded daemon")
				return
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				otelzap.S().Warn("context canceled daemon")
				return
			}
			if ctx.Err() != nil {
				otelzap.S().Error(ctx.Err())
			}
		}
	}()
	err := t.h(Context{ctx: ctx})
	if err != nil {
		otelzap.S().Errorf("Error daemon worker: %v", err)
	} else {
		otelzap.S().Infof("Success ticker daemon")
	}
}

func (t *Transport) TimeToWait() time.Duration {
	return t.Config.Timeout
}

func (t *Transport) Shutdown() error {
	t.Ticker.Stop()
	return nil
}
