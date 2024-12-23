package app

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/ihatiko/chef/core/iface"
	"github.com/ihatiko/chef/core/store"
)

type Option func(*App)

type App struct {
	context    context.Context
	Components []iface.IComponent
}
type SharedComponents map[string][]iface.IComponent

func Modules(components ...iface.IComponent) {
	app := new(App)
	app.context = context.Background()
	buffer := map[string]struct{}{}
	for _, component := range components {
		if _, ok := buffer[component.Name()]; !ok {
			buffer[component.Name()] = struct{}{}
			app.Components = append(app.Components, component)
		}
	}
	fatalState := true
	for _, pkg := range store.PackageStore.Get() {
		packageName := pkg.Name()
		if env := os.Getenv("TECH.SERVICE.DEBUG"); env != "" {
			if state, err := strconv.ParseBool(env); err == nil {
				fatalState = !state
			}
		}
		if pkg.HasError() {
			if fatalState {
				slog.Error("init package", slog.Any("error", pkg.Error()), slog.String("package", packageName))
				os.Exit(1)
			} else {
				slog.Warn("init package", slog.Any("error", pkg.Error()), slog.String("package", packageName))
			}
		}
	}
	for _, component := range app.Components {
		store.LivenessStore.Load(component)
		if component == nil {
			if fatalState {
				slog.Error("empty struct [func Deployment(components ...iface.IComponent)]")
				os.Exit(1)
			} else {
				slog.Warn("empty struct [func Deployment(components ...iface.IComponent)]")
			}
			return
		}
		go func(component iface.IComponent) {
			defer func() {
				if r := recover(); r != nil {
					if fatalState {
						slog.Error("recovered from panic", slog.Any("recover", r))
						os.Exit(1)
					} else {
						slog.Debug("recovered from panic", slog.Any("recover", r))
					}
				}
			}()
			err := component.Run()
			if err != nil {
				slog.Error("error run component", slog.Any("error", err))
				os.Exit(1)
			}
		}(component)
	}
	app.Graceful(app.Components)
}
