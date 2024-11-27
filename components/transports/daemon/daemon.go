package daemon

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const (
	defaultTimeout  = 10 * time.Second
	defaultInterval = 10 * time.Second
	defaultWorker   = 1
	componentName   = "daemon"
)

var successDaemon = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "daemon_processing_success",
}, []string{})

var failedDaemon = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "daemon_processing_failed",
}, []string{})

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
	Ticker *time.Ticker
}
type Options func(*transport)

func (t transport) Name() string {
	return fmt.Sprintf("%s id: %s", componentName, uuid.New().String())
}

func (cfg *Config) Use(opts ...Options) transport {
	t := new(transport)

	t.Config = cfg
	for _, opt := range opts {
		opt(t)
	}
	if t.Config.Timeout == 0 {
		t.Config.Timeout = defaultTimeout
	}
	if t.Config.Interval == 0 {
		t.Config.Interval = defaultInterval
	}
	t.Ticker = time.NewTicker(cfg.Interval)
	if t.Config.Workers != 0 {
		t.Config.Workers = defaultWorker
	}

	return *t
}

func (t transport) Routing(fn h) transport {
	t.h = fn
	return t
}

func (t transport) Run() {
	slog.Info("starting daemon")
	if t.h == nil {
		otelzap.L().Fatal("daemon transport handler is nil")
	}
	for range t.Ticker.C {
		wg := &sync.WaitGroup{}
		for i := range t.Config.Workers {
			go func(id int) {
				wg.Add(1)
				defer wg.Done()
				slog.Info("Start daemon worker",
					zap.Int("worker", i+1),
				)
				t.handler(i + 1)
				slog.Info("End daemon worker",
					zap.Int("worker", i+1),
				)
			}(i)
		}
		wg.Wait()
	}
}

func (t transport) handler(id int) {
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
				return
			}
			if ctx.Err() != nil {
				otelzap.S().Error(ctx.Err())
			}
		}
	}()
	err := t.h(Request{ctx: ctx, id: id})
	if err != nil {
		otelzap.S().Errorf("Error daemon worker: %v", err)
		failedDaemon.WithLabelValues().Inc()
		return
	}
	successDaemon.WithLabelValues().Inc()
}

func (t transport) Live(ctx context.Context) error {
	return nil
}

func (t transport) TimeToWait() time.Duration {
	return t.Config.Timeout
}

func (t transport) Shutdown() error {
	t.Ticker.Stop()
	return nil
}
