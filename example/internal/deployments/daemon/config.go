package daemon

import (
	"github.com/ihatiko/chef/components/clients/postgresql"
	"github.com/ihatiko/chef/components/clients/redis"
	"github.com/ihatiko/chef/components/transports/daemon"
)

type Config struct {
	Daemon          daemon.Config     `toml:"daemon"`
	ReadPostgreSQL  postgresql.Config `toml:"read-postgresql"`
	WritePostgreSQL postgresql.Config `toml:"write-postgresql"`
	Redis           redis.Config      `toml:"redis"`
}
