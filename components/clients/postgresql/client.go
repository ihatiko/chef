package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ihatiko/chef/components/core/store"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	maxOpenConnections   = 60
	connMaxLifetime      = 120
	maxIdleConnections   = 30
	connMaxIdleTime      = 20
	defaultQueryExecMode = "simple_protocol"
	defaultSslMode       = "disable"
)
const (
	postgresql    = "postgresql"
	defaultDriver = "pgx"
)

var metricProvider *metric.MeterProvider

func getMetricProvider() *metric.MeterProvider {
	if metricProvider == nil {
		exporter, err := prometheus.New()
		if err != nil {
			log.Fatal(err)
		}
		metricProvider = metric.NewMeterProvider(metric.WithReader(exporter))
	}

	return metricProvider
}

func (c *Config) toPgConnection() string {
	queryExecMode := c.QueryExecMode
	if queryExecMode == "" {
		queryExecMode = defaultQueryExecMode
	}

	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s default_query_exec_mode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Database,
		c.Password,
		c.SSLMode,
		queryExecMode,
	)
	return dataSourceName
}

type Client struct {
	Db  *sqlx.DB
	cfg *Config
	err error
}

func (c Client) Live(ctx context.Context) error {
	return c.Db.PingContext(ctx)
}
func (c Client) Error() error {
	return c.err
}
func (c Client) HasError() bool {
	return c.err != nil
}

func (c Client) Name() string {
	return fmt.Sprintf("name: %s host:%s port: %d database: %s", postgresql, c.cfg.Host, c.cfg.Port, c.cfg.Database)
}

func (c *Config) New() Client {
	pg, err := c.newConnection()
	client := Client{Db: pg, cfg: c, err: err}
	store.PackageStore.Load(client)
	return client
}

func (c *Config) newConnection() (*sqlx.DB, error) {
	if c.PgDriver == "" {
		c.PgDriver = defaultDriver
	}
	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = maxOpenConnections
	}
	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = maxIdleConnections
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxIdleTime = connMaxIdleTime * time.Second
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = connMaxLifetime * time.Second
	}
	if c.SSLMode == "" {
		c.SSLMode = defaultSslMode
	}
	connectionString := c.toPgConnection()

	db, err := otelsqlx.Connect(c.PgDriver, connectionString,
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
			attribute.KeyValue{Key: "driver", Value: attribute.StringSliceValue(sql.Drivers())},
		),
		otelsql.WithDBName(c.Database),
		otelsql.WithMeterProvider(getMetricProvider()),
	)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxIdleConnections)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	db.SetMaxIdleConns(c.MaxIdleConnections)
	db.SetConnMaxIdleTime(c.ConnMaxIdleTime)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
