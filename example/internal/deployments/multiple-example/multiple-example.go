package multiple_example

import (
	"example/config"
	"example/internal/features/peoples"
	"example/internal/features/planets"

	"github.com/ihatiko/olymp/core/app"
)

type MultipleExample struct {
	config.MultipleExample
	iPlanetsTransport planets.ITransport
	iPeoplesTransport peoples.ITransport
}

func (d MultipleExample) Run() {
	app.Modules(
		d.Daemon.Use().Routing(d.iPlanetsTransport.Load),
		d.Cron.Use().Routing(d.iPlanetsTransport.Update),
		d.PlanetsGrpcService.Use().Routing(d.iPlanetsTransport),
		d.CharactersGrpcService.Use().Routing(d.iPeoplesTransport),
	)
}
