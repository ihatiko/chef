package app

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/ihatiko/olymp/core/iface"
	"github.com/ihatiko/olymp/core/store"
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
		if env := os.Getenv("TECH_SERVICE_DEBUG"); env != "" {
			if state, err := strconv.ParseBool(env); err == nil {
				fatalState = !state
			}
		}
		if pkg.HasError() {

			if fatalState {
				slog.Error("init package", slog.Any("error", pkg.Error()), slog.Any("package", packageName))
				os.Exit(1)
			} else {
				slog.Debug("init package", slog.Any("error", pkg.Error()), slog.Any("package", packageName))
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
				slog.Debug("empty struct [func Deployment(components ...iface.IComponent)]")
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
			component.Run()
		}(component)
	}
	app.Graceful(app.Components)
}
