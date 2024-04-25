package planets

import (
	protoPlanets "example/protoc/planets"

	"github.com/ihatiko/olymp/temple/transports/cron"
	"github.com/ihatiko/olymp/temple/transports/daemon"
)

type IPlanetsTransport interface {
	protoPlanets.PlanetsServiceServer
	Load(request daemon.Request) error
	Update(request cron.Request) error
}
