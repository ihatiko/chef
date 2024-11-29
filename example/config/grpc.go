package config

import (
	protoPeoples "example/pkg/protoc/peoples"
	protoPlanets "example/pkg/protoc/planets"
	"github.com/ihatiko/olymp/components/clients/postgresql"
)

type GrpcExample struct {
	PlanetsGrpcService    protoPlanets.PlanetsConfig `toml:"grpc"`
	CharactersGrpcService protoPeoples.PeoplesConfig `toml:"grpc"`
	ReadPostgreSQL        postgresql.Config          `toml:"read-postgresql"`
	WritePostgreSQL       postgresql.Config          `toml:"write-postgresql"`
}
