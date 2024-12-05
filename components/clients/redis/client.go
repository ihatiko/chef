package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/ihatiko/olymp/core/store"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const (
	defaultReadTimeout        = 5
	defaultWriteTimeout       = 5
	defaultConnMaxLifetime    = 120
	defaultMaxIdleConnections = 30
	defaultConnMaxIdleTime    = 20
)
const (
	keyValue = "key/value"
)

type Client struct {
	Db  *redis.Client
	cfg *Config
	err error
}

func (c Client) Name() string {
	return fmt.Sprintf("name: %s host:%s sentinelAddrs: %v", keyValue, c.cfg.Addr, c.cfg.SentinelAddrs)
}

func (c Client) Live(ctx context.Context) error {
	return c.Db.Ping(ctx).Err()
}

func (c Client) Error() error {
	return c.err
}

func (c Client) HasError() bool {
	return c.err != nil
}

func (c *Config) New() Client {
	client := Client{cfg: c}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = defaultConnMaxLifetime * time.Second
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxIdleTime = defaultConnMaxIdleTime * time.Second
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout * time.Second
	}
	if c.Sentinels {
		client.Db = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:      c.MasterName,
			SentinelAddrs:   c.SentinelAddrs,
			DB:              c.Database,
			WriteTimeout:    c.WriteTimeout,
			ReadTimeout:     c.ReadTimeout,
			ConnMaxIdleTime: c.ConnMaxIdleTime,
			ConnMaxLifetime: c.ConnMaxLifetime,
		})
	} else {
		client.Db = redis.NewClient(&redis.Options{
			Addr:            c.Addr,
			Password:        c.Password,
			DB:              c.Database,
			Username:        c.UserName,
			WriteTimeout:    c.WriteTimeout,
			ReadTimeout:     c.ReadTimeout,
			ConnMaxIdleTime: c.ConnMaxIdleTime,
			ConnMaxLifetime: c.ConnMaxLifetime,
		})
	}
	if err := redisotel.InstrumentTracing(client.Db); err != nil {
		client.err = err
		return client
	}
	if err := redisotel.InstrumentMetrics(client.Db); err != nil {
		client.err = err
		return client
	}
	client.err = client.Db.Ping(context.Background()).Err()
	store.PackageStore.Load(client)
	return client
}
