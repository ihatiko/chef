package http

import "time"

type Config struct {
	Port          int           `toml:"port"`
	PprofPort     int           `toml:"pprof_port"`
	Timeout       time.Duration `toml:"timeout"`
	Pprof         bool          `toml:"pprof"`
	LivenessPath  string        `toml:"liveness_path"`
	ReadinessPath string        `toml:"readiness_path"`
	MetricsPath   string        `toml:"metrics_path"`
}
