package grpc

import (
	protoPeoples "example/pkg/protoc/peoples"
	protoPlanets "example/pkg/protoc/planets"
	"github.com/ihatiko/chef/components/clients/postgresql"
	"github.com/ihatiko/chef/components/clients/redis"
)

type Config struct {
	PlanetsGrpcService    protoPlanets.PlanetsConfig `toml:"grpc-server"`
	CharactersGrpcService protoPeoples.PeoplesConfig `toml:"grpc-server"`
	Redis                 redis.Config               `toml:"redis"`
	ReadPostgreSQL        postgresql.Config          `toml:"read-postgresql"`
	WritePostgreSQL       postgresql.Config          `toml:"write-postgresql"`
}
