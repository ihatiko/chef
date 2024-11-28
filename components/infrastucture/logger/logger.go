package logger

import (
	"log/slog"
	"os"
	"strconv"
)

type Option = func(config *Config)

type Config struct {
	Encoding string
	Level    slog.Level
}

func (c *Config) New(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	opt := new(slog.HandlerOptions)
	opt.Level = c.Level
	if env := os.Getenv("TECH_SERVICE_DEBUG"); env != "" {
		if state, err := strconv.ParseBool(env); err == nil {
			if state {
				opt.Level = slog.LevelDebug
			}
		}
	}

	h := slog.NewJSONHandler(os.Stdout, opt)
	logger := slog.New(h)
	slog.SetDefault(logger)

}
