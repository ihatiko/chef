package config

import (
	protoPeoples "example/pkg/protoc/peoples"
	protoPlanets "example/pkg/protoc/planets"
	kProducer "github.com/ihatiko/olymp/components/clients/kafka-producer"
	"github.com/ihatiko/olymp/components/clients/postgresql"
	"github.com/ihatiko/olymp/components/clients/redis"
	"github.com/ihatiko/olymp/components/transports/cron"
	"github.com/ihatiko/olymp/components/transports/daemon"
)

type MultipleExample struct {
	Cron                  cron.Config                `toml:"cron"`
	Daemon                daemon.Config              `toml:"daemon"`
	PlanetsGrpcService    protoPlanets.PlanetsConfig `toml:"grpc"`
	CharactersGrpcService protoPeoples.PeoplesConfig `toml:"grpc"`
	ReadPostgreSQL        postgresql.Config          `toml:"read-postgresql"`
	WritePostgreSQL       postgresql.Config          `toml:"write-postgresql"`
	Redis                 redis.Config               `toml:"redis"`
	PlanetsProducer       kProducer.Config           `toml:"kafka,kafka-planets-producer"`
}
