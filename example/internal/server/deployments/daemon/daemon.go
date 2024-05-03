package daemon

import (
	"example/config"
	"example/internal/features/planets"

	"github.com/ihatiko/olymp/hephaestus/app"
)

type DaemonDeploymentExample struct {
	config.DaemonDeploymentExample
	transport planets.ITransport
}

func (d DaemonDeploymentExample) Run() {
	app.Modules(
		d.Daemon.Use().Routing(d.transport.Load),
	)
}
