package config

import (
	"github.com/ihatiko/olymp/components/clients/postgresql"
	"github.com/ihatiko/olymp/components/clients/redis"
	"github.com/ihatiko/olymp/components/transports/daemon"
)

type DaemonDeploymentExample struct {
	Daemon          daemon.Config     `toml:"daemon"`
	ReadPostgreSQL  postgresql.Config `toml:"read-postgresql"`
	WritePostgreSQL postgresql.Config `toml:"write-postgresql"`
	Redis           redis.Config      `toml:"redis"`
}
