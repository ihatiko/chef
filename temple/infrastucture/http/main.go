package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/ihatiko/olymp/hephaestus/iface"
	store "github.com/ihatiko/olymp/hephaestus/store"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
)

const (
	metrics   = "/metrics"
	live      = "/liveness"
	readiness = "/readiness"
)
const (
	defaultPprofPort = 8081
	defaultPort      = 8080
)

const (
	defaultTimeout = time.Second * 15
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

type Live struct {
	InternalError error `json:"internal_error"`
	ContextError  error `json:"context_error"`
}

func (t *Transport) Live(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	mutex := &sync.Mutex{}
	defer mutex.Unlock()
	result := map[string]Live{}
	status := true
	for _, v := range store.LivenessStore.Get() {
		wg.Add(1)
		go func(iLive iface.ILive) {
			resultLv := Live{}
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.TODO(), t.Config.Timeout)
			defer cancel()
			go func() {
				select {
				case <-ctx.Done():
					if errors.Is(ctx.Err(), context.DeadlineExceeded) {
						resultLv.InternalError = fmt.Errorf("context context deadline exceeded tech-http component: %s", iLive.Name())
						return
					}
					if errors.Is(ctx.Err(), context.Canceled) {
						return
					}
					if ctx.Err() != nil {
						fmt.Errorf("context errored tech-http component: %s", iLive.Name())
					}
				}
			}()
			resultLv.InternalError = iLive.Live(ctx)
			if resultLv.InternalError != nil || resultLv.ContextError != nil {
				status = false
			}
			mutex.Lock()
			defer mutex.Unlock()
			result[iLive.Name()] = resultLv
		}(v)
	}
	wg.Wait()
	data, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if !status {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (t *Transport) Run() {
	if t.Config.Pprof {
		go func() {
			pprofMux := http.NewServeMux()
			if t.Config.PprofPort == 0 {
				t.Config.PprofPort = defaultPprofPort
			}
			fmt.Println(fmt.Sprintf("Start pprof server port: %d", t.Config.PprofPort))
			pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
			pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
			pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", t.Config.PprofPort), pprofMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Println(fmt.Sprintf("close http server %v", err))
			}
		}()
	}
	go func() {
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
		if t.Config.Timeout == 0 {
			t.Config.Timeout = defaultTimeout
		}
		mux.Handle(metricsPath, promhttp.Handler())
		mux.HandleFunc(readinessPath, t.Ready)
		mux.HandleFunc(livenessPath, t.Live)
		fmt.Println(fmt.Sprintf("Start http server port: %d", t.Config.Port))
		if t.Config.Port == 0 {
			t.Config.Port = defaultPort
		}
		port := fmt.Sprintf(":%d", t.Config.Port)
		server := http.Server{
			Addr:         port,
			Handler:      mux,
			ReadTimeout:  t.Config.Timeout,
			WriteTimeout: t.Config.Timeout,
		}
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(fmt.Sprintf("close http server %v", err))
		}
	}()
}
