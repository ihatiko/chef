package config

import (
	protoCharacters "example/protoc/characters"
	protoPlanets "example/protoc/planets"

	"github.com/ihatiko/olymp/temple/infrastucture/postgresql"
	"github.com/ihatiko/olymp/temple/transports/cron"
	"github.com/ihatiko/olymp/temple/transports/daemon"
)

type MultipleExample struct {
	Cron                  cron.Config                      `toml:"cron"`
	Daemon                daemon.Config                    `toml:"daemon"`
	PlanetsGrpcService    protoPlanets.PlanetsConfig       `toml:"grpc"`
	CharactersGrpcService protoCharacters.CharactersConfig `toml:"grpc"`
	ReadPostgreSql        postgresql.Config                `toml:"read-postgresql"`
	WritePostgreSql       postgresql.Config                `toml:"read-postgresql"`
}
