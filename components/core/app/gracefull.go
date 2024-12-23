package app

import (
	"go.uber.org/zap"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ihatiko/chef/core/iface"
)

func (a *App) Graceful(components []iface.IComponent) {
	<-a.Wait()
	a.BeforeShutdown(components)
	slog.Info("starting shutdown ...")
	a.Shutdown(components)

	slog.Info("starting delay [terminating old requests] ...")
	Delay(
		components...,
	)
	slog.Info("starting delay [terminating old requests] ... done")
	a.AfterShutdown(components)
	slog.Info("starting shutdown ... done")
	slog.Info("Server exit properly")
}

func (a *App) AfterShutdown(components []iface.IComponent) {
	for _, t := range components {
		if component, ok := t.(iface.IAfterLifecycleComponent); ok {
			componentName := component.Name()
			slog.Info("starting after shutdown...", slog.String("component", componentName))
			err := component.AfterShutdown()
			if err != nil {
				slog.Error("shutdown", slog.Any("error", err), slog.String("component", componentName))
			}
			slog.Info("starting after shutdown...done", slog.String("component", componentName))
		}
	}
}

func (a *App) Shutdown(components []iface.IComponent) {
	for _, component := range components {
		componentName := component.Name()
		slog.Info("starting shutdown...", slog.String("component", componentName))
		err := component.Shutdown()
		if err != nil {
			slog.Error("shutdown", slog.Any("error", err), slog.String("component", componentName))
		}
		slog.Info("starting shutdown...done", slog.String("component", componentName))
	}
}

func (a *App) BeforeShutdown(components []iface.IComponent) {
	for _, t := range components {
		if component, ok := t.(iface.IBeforeLifecycleComponent); ok {
			slog.Info("starting before shutdown...", slog.String("component", component.Name()))
			err := component.BeforeShutdown()
			if err != nil {
				slog.Error("error before shutdown", zap.Any("error", err), slog.String("component", component.Name()))
			}
			slog.Info("starting before shutdown...done", slog.String("component", component.Name()))
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
