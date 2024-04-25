package app

import (
	"context"

	"github.com/ihatiko/olymp/hephaestus/iface"
	"github.com/ihatiko/olymp/hephaestus/store"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type Option func(*App)

type App struct {
	context    context.Context
	Components []iface.IComponent
}

func Deployment(components ...iface.IComponent) {
	app := new(App)
	app.context = context.Background()

	for _, o := range app.Components {
		store.LivenessStore.Load(o)
		if o == nil {
			otelzap.L().Fatal("empty struct [func Deployment(components ...iface.IComponent)]")
			return
		}
		go func(component iface.IComponent) {
			o.Run()
		}(o)
	}
	app.Graceful(app.Components)
}
