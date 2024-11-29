package grpc

import (
	"example/config"
	"example/internal/features/peoples"
	"example/internal/features/planets"
	"github.com/ihatiko/olymp/core/app"
)

type GrpcExample struct {
	config.GrpcExample
	iPlanetsTransport planets.ITransport
	iPeoplesTransport peoples.ITransport
}

func (d GrpcExample) Run() {
	app.Modules(
		d.PlanetsGrpcService.Use().Routing(d.iPlanetsTransport),
		d.CharactersGrpcService.Use().Routing(d.iPeoplesTransport),
	)
}
