package daemon

import (
	"example/config"
	"example/internal/features/planets"

	"github.com/ihatiko/olymp/hephaestus/app"
)

type MultipleExample struct {
	config.MultipleExample
	iplanetsTransport planets.IPlanetsTransport
}

func (d MultipleExample) Run() {
	app.Deployment(
		d.Daemon.Use().Routing(d.iplanetsTransport.Load),
		d.Cron.Use().Routing(d.iplanetsTransport.Update),
		d.PlanetsGrpc.Use().Routing(d.iplanetsTransport),
	)
}
