package app

import (
	"context"

	"github.com/ihatiko/olymp/hephaestus/iface"
	"github.com/ihatiko/olymp/hephaestus/store"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
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
	for _, component := range app.Components {
		store.LivenessStore.Load(component)
		if component == nil {
			otelzap.L().Fatal("empty struct [func Deployment(components ...iface.IComponent)]")
			return
		}
		go func(component iface.IComponent) {
			defer func() {
				if r := recover(); r != nil {
					otelzap.L().Error("recovered from panic", zap.Any("recover", r))
				}
			}()
			component.Run()
		}(component)
	}
	app.Graceful(app.Components)
}