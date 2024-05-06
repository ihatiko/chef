package config

import "github.com/ihatiko/olymp/components/transports/cron"

type Cron struct {
	Cron cron.Config `toml:"cron"`
}
