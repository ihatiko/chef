package daemon

import (
	"context"
	"errors"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const (
	defaultInterval = 10 * time.Second
)

type Config struct {
	Timeout  time.Duration
	Interval time.Duration
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

	return t
}

func Router[T any](router H) Options {
	return func(t *Transport) {
		t.h = router
	}
}

func (t *Transport) Run() {
	for range t.Ticker.C {
		otelzap.S().Infof("Start ticker worker")
		t.handler()
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
			if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
				otelzap.S().Error(ctx.Err())
			}
		}
	}()
	err := t.h(Context{ctx: ctx})
	if err != nil {
		otelzap.S().Errorf("Error ticker worker: %v", err)
	} else {
		otelzap.S().Infof("Success ticker worker")
	}
}

func (t *Transport) TimeToWait() time.Duration {
	return time.Duration(t.Config.Timeout) * time.Second
}

func (t *Transport) Shutdown() error {
	t.Ticker.Stop()
	return nil
}
