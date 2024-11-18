package daemon

import (
	"example/config"
	"example/internal/features/planets"
	"github.com/ihatiko/olymp/core/app"
)

type Daemon struct {
	config.DaemonDeploymentExample
	transport planets.ITransport
}

func (d Daemon) Run() {
	app.Modules(
		d.Daemon.Use().Routing(d.transport.Load),
	)
}
