package config

import (
	protoPeoples "example/pkg/protoc/peoples"
	protoPlanets "example/pkg/protoc/planets"
	"github.com/ihatiko/olymp/components/clients/postgresql"
)

type Grpc struct {
	PlanetsGrpcService    protoPlanets.PlanetsConfig `toml:"grpc-server"`
	CharactersGrpcService protoPeoples.PeoplesConfig `toml:"grpc-server"`
	ReadPostgreSQL        postgresql.Config          `toml:"read-postgresql"`
	WritePostgreSQL       postgresql.Config          `toml:"write-postgresql"`
}
