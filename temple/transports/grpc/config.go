package grpc

import "time"

type Metrics struct {
	EnableHandlingTimeHistogram       bool
	EnableClientHandlingTimeHistogram bool
}

type Config struct {
	Port              string
	TimeOut           time.Duration
	Time              time.Duration
	MaxConnectionAge  time.Duration
	MaxConnectionIdle time.Duration
	Healthz           bool
	Reflect           bool
	MaxRecvMsgSize    int
	Metrics           Metrics
	Shared            *bool
}
