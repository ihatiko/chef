package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ihatiko/olymp/core/iface"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func (a *App) Graceful(components []iface.IComponent) {
	<-a.Wait()
	slog.Info("starting graceful before components...")
	a.BeforeShutdown(components)
	slog.Info("starting graceful before components... done")
	slog.Info("starting shutdown ...")
	a.Shutdown(components)
	slog.Info("starting shutdown ... done")
	slog.Info("starting delay [terminating old requests] ...")
	Delay(
		components...,
	)
	slog.Info("starting delay [terminating old requests] ... done")
	slog.Info("starting gracefull after components...")
	a.AfterShutdown(components)
	slog.Info("starting gracefull after components... done")
	slog.Info("Server exit properly")
}

func (a *App) AfterShutdown(components []iface.IComponent) {
	for _, t := range components {
		if component, ok := t.(iface.IAfterLifecycleComponent); ok {
			slog.Info("starting after shutdown...", zap.String("component", component.Name()))
			err := component.AfterShutdown()
			if err != nil {
				otelzap.S().Error(err)
			}
			slog.Info("starting after shutdown...done", zap.String("component", component.Name()))
		}
	}
}

func (a *App) Shutdown(components []iface.IComponent) {
	for _, component := range components {
		slog.Info("starting shutdown...", zap.String("component", component.Name()))
		err := component.Shutdown()
		if err != nil {
			otelzap.S().Error(err)
		}
		slog.Info("starting shutdown...done", zap.String("component", component.Name()))
	}
}

func (a *App) BeforeShutdown(components []iface.IComponent) {
	for _, t := range components {
		if component, ok := t.(iface.IBeforeLifecycleComponent); ok {
			slog.Info("starting before shutdown...", zap.String("component", component.Name()))
			err := component.BeforeShutdown()
			if err != nil {
				otelzap.S().Error(err)
			}
			slog.Info("starting before shutdown...done", zap.String("component", component.Name()))
		}
	}
}
func (a *App) Wait() chan struct{} {
	result := make(chan struct{}, 1)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		<-quit
		result <- struct{}{}
	}()
	go func() {
		<-a.context.Done()
		result <- struct{}{}
	}()
	return result
}

func Delay(times ...iface.IComponent) {
	var cur time.Duration
	for _, dur := range times {
		d := dur.TimeToWait()
		if d > cur {
			cur = d
		}
	}
	time.Sleep(cur)
}
