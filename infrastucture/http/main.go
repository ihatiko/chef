package open_api_transport

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metrics   = "/metrics"
	live      = "/liveness"
	readiness = "/readiness"
)
const (
	defaultPprofPort = ":8080"
	defaultPort      = ":8081"
)

type Transport struct {
	Config *Config
}
type Options func(*Transport)

func (cfg *Config) Use(opts ...Options) *Transport {
	t := &Transport{
		Config: cfg,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *Transport) Ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (t *Transport) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (t *Transport) Run() {
	mux := http.NewServeMux()
	metricsPath := metrics
	if t.Config.MetricsPath != "" {
		metricsPath = t.Config.MetricsPath
	}
	livenessPath := live
	if t.Config.LivenessPath != "" {
		livenessPath = t.Config.LivenessPath
	}
	readinessPath := readiness
	if t.Config.ReadinessPath != "" {
		readinessPath = t.Config.ReadinessPath
	}
	mux.Handle(metricsPath, promhttp.Handler())
	mux.HandleFunc(readinessPath, t.Ready)
	mux.HandleFunc(livenessPath, t.Live)
	fmt.Println(fmt.Sprintf("Start http server %s", t.Config.Port))
	if t.Config.Port == "" {
		t.Config.Port = defaultPort
	}
	if t.Config.Pprof {
		go func() {
			pprofMux := http.NewServeMux()
			if t.Config.PprofPort == "" {
				t.Config.PprofPort = defaultPprofPort
			}
			pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
			pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
			if err := http.ListenAndServe(t.Config.PprofPort, pprofMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Println(fmt.Sprintf("close http server %v", err))
			}
		}()
	}
	if err := http.ListenAndServe(t.Config.Port, mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println(fmt.Sprintf("close http server %v", err))
	}
}
