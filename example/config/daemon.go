package config

import (
	"github.com/ihatiko/olymp/temple/infrastucture/postgresql"
	"github.com/ihatiko/olymp/temple/transports/daemon"
)

type DaemonDeploymentExample struct {
	Daemon          daemon.Config     `toml:"daemon"`
	ReadPostgreSQL  postgresql.Config `toml:"read-postgresql"`
	WritePostgreSQL postgresql.Config `toml:"write-postgresql"`
}
