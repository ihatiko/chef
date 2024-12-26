package grpc

import (
	"example/internal/features/peoples"
	"example/internal/features/planets"
	"github.com/ihatiko/chef/components/core/app"
)

type Deployment struct {
	Config
	iPlanetsTransport planets.ITransport
	iPeoplesTransport peoples.ITransport
}

func (d Deployment) Run() {
	app.Modules(
		d.PlanetsGrpcService.Setup(d.iPlanetsTransport),
		d.CharactersGrpcService.Setup(d.iPeoplesTransport),
	)
}
