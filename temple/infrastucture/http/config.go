package http

import "time"

type Config struct {
	Port          string
	PprofPort     string
	Timeout       time.Duration
	Pprof         bool
	LivenessPath  string
	ReadinessPath string
	MetricsPath   string
}
