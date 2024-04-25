package daemon

import "time"

type Config struct {
	Timeout  time.Duration
	Interval time.Duration
	Workers  int
}
