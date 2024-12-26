package daemon

import (
	"example/internal/features/planets"
	"github.com/ihatiko/chef/components/core/app"
)

type Deployment struct {
	Config
	transport planets.ITransport
}

func (d Deployment) Run() {
	app.Modules(
		d.Daemon.Setup(d.transport.Load),
	)
}
