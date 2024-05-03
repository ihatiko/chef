package cron

import "time"

type Config struct {
	Timeout time.Duration `toml:"timeout"`
	Workers int           `toml:"workers"`
}
