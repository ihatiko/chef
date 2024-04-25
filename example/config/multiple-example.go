package config

import (
	protoPlanets "example/protoc/planets"

	"github.com/ihatiko/olymp/temple/transports/cron"
	"github.com/ihatiko/olymp/temple/transports/daemon"
	"github.com/ihatiko/olymp/temple/transports/grpc"
)

type MultipleExample struct {
	Cron               cron.Config
	Daemon             daemon.Config
	Grpc               grpc.Config
	PlanetsGrpcService protoPlanets.PlanetsConfig
}
