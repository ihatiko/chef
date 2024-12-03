package grpc

import "time"

type Config struct {
	Port                string        `toml:"port"`
	TimeOut             time.Duration `toml:"timeout"`
	TimeKeepaliveParams time.Duration `toml:"time"`
	MaxConnectionAge    time.Duration `toml:"max_connection_age"`
	MaxConnectionIdle   time.Duration `toml:"max_connection_idle"`
	Healthz             bool          `toml:"healthz"`
	Reflect             *bool         `toml:"reflect"`
	MaxRecMessageSize   int           `toml:"max_rec_message_size"`
	Metrics             struct {
		EnableHandlingTimeHistogram       bool `toml:"enable_handling_time_histogram"`
		EnableClientHandlingTimeHistogram bool `toml:"enable_client_handling_time_histogram"`
	} `toml:"metrics"`
}

func (c *Config) IsValid() bool {
	if c.Port == "" {
		return false
	}
	return true
}
