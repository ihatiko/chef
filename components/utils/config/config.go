package config

import (
	"github.com/ihatiko/components/utils/environ"
	"os"
	"path"

	"github.com/ihatiko/olymp/infrastucture/components/utils/toml"
)

const (
	configPath = "config.toml"
)

type Settings struct {
	Path string
}

type Options func(setting *Settings)

func WithPath(path string) Options {
	return func(setting *Settings) {
		setting.Path = path
	}
}

func ToConfig[T any](t T, opts ...Options) error {
	s := new(Settings)
	for _, opt := range opts {
		opt(s)
	}
	if s.Path == "" {
		p, err := os.Getwd()
		if err != nil {
			return err
		}
		// TODO support various OS
		s.Path = path.Join(p, configPath)
	}
	f, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(f, t)
	if err != nil {
		return err
	}
	environ.Parse(t)
	return err
}
