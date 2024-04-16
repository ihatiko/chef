package http

type Config struct {
	Port          string
	PprofPort     string
	Pprof         bool
	LivenessPath  string
	ReadinessPath string
	MetricsPath   string
}
