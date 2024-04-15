package open_api_transport

type Config struct {
	Port          string
	PprofPort     string
	Pprof         bool
	LivenessPath  string
	ReadinessPath string
	MetricsPath   string
}
