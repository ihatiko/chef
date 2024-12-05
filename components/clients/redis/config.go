package redis

import "time"

type Config struct {
	Host            string        `json:"addr"`
	Login           string        `toml:"user_name"`
	Password        string        `toml:"password"`
	Database        int           `toml:"database"`
	SentinelAddrs   []string      `toml:"sentinel_addrs"`
	MasterName      string        `toml:"master_name"`
	Sentinels       bool          `toml:"sentinels"`
	DialTimeout     time.Duration `toml:"dial_timeout"`
	ReadTimeout     time.Duration `toml:"read_timeout"`
	WriteTimeout    time.Duration `toml:"write_timeout"`
	ConnMaxIdleTime time.Duration `toml:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"`
	MaxIdleConns    int           `toml:"max_idle_conns"`
}
