package daemon

import "time"

type Config struct {
	Timeout  time.Duration `toml:"timeout"`
	Interval time.Duration `toml:"interval"`
	Workers  int           `toml:"workers"`
}
