package daemon

import (
	"example/internal/features/planets"
	"github.com/ihatiko/olymp/core/app"
)

type Deployment struct {
	Config
	transport planets.ITransport
}

func (d Deployment) Run() {
	app.Modules(
		d.Daemon.Use().Routing(d.transport.Load),
	)
}
