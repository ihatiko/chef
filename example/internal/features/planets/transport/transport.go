package transport

import (
	"context"
	"example/internal/features/planets"
	protoPlanets "example/pkg/protoc/planets"
	"fmt"

	"github.com/ihatiko/chef/components/transports/cron"
	"github.com/ihatiko/chef/components/transports/daemon"
)

type transport struct {
	service planets.IService
}

func New(service planets.IService) planets.ITransport {
	return &transport{
		service: service,
	}
}

func (t transport) UpdatePlanet(ctx context.Context, request *protoPlanets.UpdatePlanetRequest) (*protoPlanets.UpdatePlanetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (t transport) Load(request daemon.Request) error {
	fmt.Println(123)
	return nil
}

func (t transport) Update(request cron.Request) error {
	fmt.Println(123)
	return nil
}
