package config

import (
	protoPeoples "example/pkg/protoc/peoples"
	protoPlanets "example/pkg/protoc/planets"
	"github.com/ihatiko/olymp/components/clients/postgresql"
	"github.com/ihatiko/olymp/components/clients/redis"
)

type Grpc struct {
	PlanetsGrpcService    protoPlanets.PlanetsConfig `toml:"grpc-server"`
	CharactersGrpcService protoPeoples.PeoplesConfig `toml:"grpc-server"`
	Redis                 redis.Config               `toml:"redis"`
	ReadPostgreSQL        postgresql.Config          `toml:"read-postgresql"`
	WritePostgreSQL       postgresql.Config          `toml:"write-postgresql"`
}
