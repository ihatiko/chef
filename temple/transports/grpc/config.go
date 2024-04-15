package grpc_transport

import "time"

type Config struct {
	Port              string
	TimeOut           time.Duration
	MaxConnectionAge  time.Duration
	MaxConnectionIdle time.Duration
}
