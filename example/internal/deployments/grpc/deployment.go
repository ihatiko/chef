package grpc

import (
	"example/internal/features/peoples"
	"example/internal/features/planets"
	"github.com/ihatiko/olymp/core/app"
)

type Deployment struct {
	Config
	iPlanetsTransport planets.ITransport
	iPeoplesTransport peoples.ITransport
}

func (d Deployment) Run() {
	app.Modules(
		d.PlanetsGrpcService.Use().Routing(d.iPlanetsTransport),
		d.CharactersGrpcService.Use().Routing(d.iPeoplesTransport),
	)
}
