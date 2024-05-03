package transport

import (
	"context"
	"example/internal/features/planets"
	protoPlanets "example/protoc/planets"

	"github.com/ihatiko/olymp/temple/transports/cron"
	"github.com/ihatiko/olymp/temple/transports/daemon"
)

type Transport struct {
	service planets.IService
}

func New(service planets.IService) planets.ITransport {
	return &Transport{}
}

func (t Transport) UpdatePlanet(ctx context.Context, request *protoPlanets.UpdatePlanetRequest) (*protoPlanets.UpdatePlanetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t Transport) Load(request daemon.Request) error {
	//TODO implement me
	panic("implement me")
}

func (t Transport) Update(request cron.Request) error {
	//TODO implement me
	panic("implement me")
}
