package planets

import (
	protoPlanets "example/pkg/protoc/planets"
	"github.com/ihatiko/olymp/components/transports/cron"
	"github.com/ihatiko/olymp/components/transports/daemon"
)

type ITransport interface {
	protoPlanets.PlanetsServiceServer
	Load(request daemon.Request) error
	Update(request cron.Request) error
}
