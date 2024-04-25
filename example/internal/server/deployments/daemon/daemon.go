package daemon

import (
	"example/config"
	"example/internal/features/planets"

	"github.com/ihatiko/olymp/hephaestus/app"
)

type DaemonDeploymentExample struct {
	config.DaemonDeploymentExample
	transport planets.IPlanetsTransport
}

func (d DaemonDeploymentExample) Run() {
	app.Deployment(
		d.Daemon.Use().Routing(d.transport.Load),
	)
}
