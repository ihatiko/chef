package planets

import (
	protoPlanets "example/pkg/protoc/planets"
	"github.com/ihatiko/chef/components/transports/cron"
	"github.com/ihatiko/chef/components/transports/daemon"
)

type ITransport interface {
	protoPlanets.PlanetsServiceServer
	Load(request daemon.Request) error
	Update(request cron.Request) error
}
