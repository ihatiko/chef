package logger

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Option = func(config *Config)

type Config struct {
	Encoding string     `toml:"encoding"`
	Level    slog.Level `toml:"level"`
}

func (c *Config) New(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	opt := new(slog.HandlerOptions)
	opt.Level = c.Level
	opt.AddSource = false
	if env := os.Getenv("TECH.SERVICE.DEBUG"); env != "" {
		if state, err := strconv.ParseBool(env); err == nil {
			if state {
				opt.Level = slog.LevelDebug
			}
		}
	}

	if strings.ToLower(c.Encoding) == "json" {
		h := slog.NewJSONHandler(os.Stdout, opt)
		logger := slog.New(h)
		slog.SetDefault(logger)
	} else {
		h := tint.NewHandler(os.Stdout, &tint.Options{
			Level:     opt.Level,
			AddSource: false,
		})
		logger := slog.New(h)
		slog.SetDefault(logger)
	}
}
