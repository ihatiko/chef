package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/google/uuid"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogger "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	healthz "github.com/ihatiko/olymp/components/transports/grpc/protoc/health"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const componentName = "grpc"

var transport map[string]Transport = make(map[string]Transport)
var active map[string]struct{} = make(map[string]struct{})
var shutdown map[string]struct{} = make(map[string]struct{})

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
	otelzap.L().Error(
		"panic occurred",
		zap.String("panic_id", panicID),
		zap.Any("panic", p),
		zap.ByteString("stacktrace", debug.Stack()),
	)
	return status.Errorf(codes.Internal, "panic (id: %s)", panicID)
}
func (t Transport) Name() string {
	return fmt.Sprintf("%s port: %s", componentName, t.Cfg.Port)
}

// Инициализация транспортного слоя grpc
func (cfg *Config) Use(
	opts ...Options,
) Transport {
	mt.Lock()
	if cfg.Port == "" {
		cfg.Port = defaultPort
	}
	if t, ok := transport[cfg.Port]; ok {
		defer mt.Unlock()
		return t
	}
	t := new(Transport)
	t.Cfg = cfg
	if t.Cfg.Metrics.EnableHandlingTimeHistogram {
		grpcPrometheus.EnableHandlingTimeHistogram()
	}
	if t.Cfg.Metrics.EnableClientHandlingTimeHistogram {
		grpcPrometheus.EnableClientHandlingTimeHistogram()
	}

	t.App = grpc.NewServer(
		grpc.MaxRecvMsgSize(t.Cfg.MaxRecvMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: t.Cfg.MaxConnectionIdle,
			Timeout:           t.Cfg.TimeOut,
			MaxConnectionAge:  t.Cfg.MaxConnectionAge,
			Time:              t.Cfg.Time,
		}),
		grpc.ChainStreamInterceptor(
			grpcRecovery.StreamServerInterceptor(
				grpcRecovery.WithRecoveryHandler(logPanic),
			),
			grpcMiddleware.ChainStreamServer(),
			grpcOpentracing.StreamServerInterceptor(),
			grpcCtxtags.StreamServerInterceptor(),
			grpcPrometheus.StreamServerInterceptor,
			otelgrpc.StreamServerInterceptor(),
			grpcLogger.StreamServerInterceptor(
				otelzap.L().Logger,
			),
		),
		grpc.ChainUnaryInterceptor(
			grpcRecovery.UnaryServerInterceptor(
				grpcRecovery.WithRecoveryHandler(logPanic),
			),
			grpcCtxtags.UnaryServerInterceptor(),
			grpcOpentracing.UnaryServerInterceptor(),
			grpcPrometheus.UnaryServerInterceptor,
			otelgrpc.UnaryServerInterceptor(),
			grpcLogger.UnaryServerInterceptor(
				otelzap.L().Logger,
			),
		),
	)

	for _, rt := range opts {
		rt(t)
	}
	if t.Cfg.Reflect {
		reflection.Register(t.App)
	}
	transport[cfg.Port] = *t
	mt.Unlock()
	return *t
}

// Native
func (t Transport) Routing(registrar grpc.ServiceDesc, impl any) Transport {
	t.App.RegisterService(&registrar, impl)
	return t
}

// Запуск транспортного слоя
func (t Transport) Run() {
	mt.Lock()
	if _, ok := active[t.Cfg.Port]; ok {
		defer mt.Unlock()
		return
	}
	otelzap.S().Infof("starting gRPC Transport port: %s ...", t.Cfg.Port)
	listener, err := net.Listen("tcp", t.Cfg.Port)
	if err != nil {
		otelzap.S().Fatal(err)
	}
	if t.Cfg.Healthz {
		otelzap.S().Info("init healthz check point ...")
		d := new(healthz.Health)
		healthz.RegisterHealthServer(t.App, d)
		otelzap.S().Info("init healthz check point ... done")
	}
	active[t.Cfg.Port] = struct{}{}
	mt.Unlock()
	err = t.App.Serve(listener)
	if err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			return
		}
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			otelzap.S().Warn(err)
			return
		}
		otelzap.S().Fatal(err)
	}
}

// Ожидание сколько нужно ждать перед выключением сервера
func (t Transport) TimeToWait() time.Duration {
	return t.Cfg.TimeOut * time.Second
}

// Безопасное выключение сервера (gracefull)
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
