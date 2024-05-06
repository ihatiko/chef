package http

import "time"

type Config struct {
	Port          int
	PprofPort     int
	Timeout       time.Duration
	Pprof         bool
	LivenessPath  string
	ReadinessPath string
	MetricsPath   string
}
