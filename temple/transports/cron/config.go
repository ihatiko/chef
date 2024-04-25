package cron

import "time"

type Config struct {
	Timeout time.Duration
	Workers int
}
