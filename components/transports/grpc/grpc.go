package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/google/uuid"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	healthz "github.com/ihatiko/olymp/components/transports/grpc/protoc/health"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const componentName = "grpc-server"

var transport = make(map[string]Transport)
var active = make(map[string]struct{})
var shutdown = make(map[string]struct{})

var mt sync.Mutex

type Transport struct {
	App *grpc.Server
	Cfg *Config
}

func (t Transport) Live(ctx context.Context) error {
	grpcClient, err := grpc.NewClient(t.Cfg.Port)
	if err != nil {
		return err
	}
	healthZClient := healthz.NewHealthClient(grpcClient)
	resp, err := healthZClient.Check(ctx, new(healthz.HealthCheckRequest))
	if err != nil {
		return err
	}
	if resp.Status != healthz.HealthCheckResponse_SERVING {
		return fmt.Errorf("health check failed: %v", resp.Status)
	}
	return err
}

type Options func(*Transport)

func logPanic(p any) error {
	panicID := uuid.New().String()
	slog.Error(
		"panic occurred",
		slog.String("panic", panicID),
		slog.Any("panic", p),
		slog.String("stack", string(debug.Stack())),
	)
	return status.Errorf(codes.Internal, "panic (id: %s)", panicID)
}
func (t Transport) Name() string {
	return fmt.Sprintf("%s port: %s", componentName, t.Cfg.Port)
}

var (
	defaultMaxConnectionAge  = 10 * time.Second
	defaultMaxConnectionIdle = 10 * time.Second
	defaultKeepaliveParams   = 10 * time.Second
)

// Use Инициализация транспортного слоя grpc
func (c *Config) Use(
	opts ...Options,
) Transport {
	mt.Lock()
	if c.Port == "" {
		c.Port = defaultPort
	}
	if c.Reflect == nil {
		reflectState := true
		c.Reflect = &reflectState
	}
	if c.MaxConnectionIdle == 0 {
		c.MaxConnectionIdle = defaultMaxConnectionIdle
	}
	if c.MaxConnectionAge == 0 {
		c.MaxConnectionAge = defaultMaxConnectionAge
	}
	if c.TimeKeepaliveParams == 0 {
		c.TimeKeepaliveParams = defaultKeepaliveParams
	}
	if t, ok := transport[c.Port]; ok {
		defer mt.Unlock()
		return t
	}
	t := new(Transport)
	t.Cfg = c
	if t.Cfg.Metrics.EnableHandlingTimeHistogram {
		grpcPrometheus.EnableHandlingTimeHistogram()
	}
	if t.Cfg.Metrics.EnableClientHandlingTimeHistogram {
		grpcPrometheus.EnableClientHandlingTimeHistogram()
	}

	t.App = grpc.NewServer(
		grpc.MaxRecvMsgSize(t.Cfg.MaxRecMessageSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: t.Cfg.MaxConnectionIdle,
			Timeout:           t.Cfg.TimeOut,
			MaxConnectionAge:  t.Cfg.MaxConnectionAge,
			Time:              t.Cfg.TimeKeepaliveParams,
		}),
		grpc.ChainStreamInterceptor(
			grpcRecovery.StreamServerInterceptor(
				grpcRecovery.WithRecoveryHandler(logPanic),
			),
			grpcMiddleware.ChainStreamServer(),
			grpcOpentracing.StreamServerInterceptor(),
			grpcCtxtags.StreamServerInterceptor(),
			grpcPrometheus.StreamServerInterceptor,
			// TODO logger
		),
		grpc.ChainUnaryInterceptor(
			grpcRecovery.UnaryServerInterceptor(
				grpcRecovery.WithRecoveryHandler(logPanic),
			),
			grpcCtxtags.UnaryServerInterceptor(),
			grpcOpentracing.UnaryServerInterceptor(),
			grpcPrometheus.UnaryServerInterceptor,
			//TODO logger
		),
	)

	for _, rt := range opts {
		rt(t)
	}
	if t.Cfg.Reflect != nil && *t.Cfg.Reflect {
		reflection.Register(t.App)
	}
	transport[c.Port] = *t
	mt.Unlock()
	return *t
}

// Routing Native
func (t Transport) Routing(registrar grpc.ServiceDesc, impl any) Transport {
	t.App.RegisterService(&registrar, impl)
	return t
}

// Run Запуск транспортного слоя
func (t Transport) Run() error {
	mt.Lock()
	if _, ok := active[t.Cfg.Port]; ok {
		defer mt.Unlock()
		return nil
	}
	slog.Info("Starting gRPC Transport...", slog.String("port", t.Cfg.Port))
	listener, err := net.Listen("tcp", t.Cfg.Port)
	if err != nil {
		slog.Error("failed to listen", slog.String("port", t.Cfg.Port), slog.Any("error", err))
		os.Exit(1)
	}
	if t.Cfg.Healthz {
		slog.Info("init healthz check point ...")
		d := new(healthz.Health)
		healthz.RegisterHealthServer(t.App, d)
		slog.Info("init healthz check point ... done")
	}
	active[t.Cfg.Port] = struct{}{}
	mt.Unlock()
	err = t.App.Serve(listener)
	if err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("gRPC Transport is stopped error: %v", err)
		}
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("gRPC Transport is stopped error: %v", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("gRPC Transport is stopped error: %v", err)
		}
		return fmt.Errorf("gRPC Transport failed: %v", err)
	}
	return nil
}

// TimeToWait Ожидание сколько нужно ждать перед выключением сервера
func (t Transport) TimeToWait() time.Duration {
	return t.Cfg.TimeKeepaliveParams
}

// Shutdown Безопасное выключение сервера (graceful)
func (t Transport) Shutdown() error {
	mt.Lock()
	if _, ok := shutdown[t.Cfg.Port]; ok {
		defer mt.Unlock()
		return nil
	}
	t.App.GracefulStop()
	shutdown[t.Cfg.Port] = struct{}{}
	mt.Unlock()
	return nil
}
