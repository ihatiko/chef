package multiple_example

import (
	"example/config"
	"example/internal/features/characters"
	"example/internal/features/planets"

	"github.com/ihatiko/olymp/core/app"
)

type MultipleExample struct {
	config.MultipleExample
	iplanetsTransport    planets.ITransport
	icharactersTransport characters.ICharactersTransport
}

func (d MultipleExample) Run() {
	app.Modules(
		d.Daemon.Use().Routing(d.iplanetsTransport.Load),
		d.Cron.Use().Routing(d.iplanetsTransport.Update),
		d.PlanetsGrpcService.Use().Routing(d.iplanetsTransport),
		d.CharactersGrpcService.Use().Routing(d.icharactersTransport),
	)
}
