package config

import (
	"errors"
	"net/url"
	"os"

	toml "github.com/pelletier/go-toml"
)

const (
	configPath = "config/config.toml"
)

type ConfigSettings struct {
	Path string
}

type Options func(setting *ConfigSettings)

func WithPath(path string) Options {
	return func(setting *ConfigSettings) {
		setting.Path = path
	}
}

func ToConfig[T any](t T, opts ...Options) error {
	s := new(ConfigSettings)
	for _, opt := range opts {
		opt(s)
	}
	if s.Path == "" {
		path, err := os.Getwd()
		if err != nil {
			return err
		}
		newPath, err := url.JoinPath(path, configPath)
		if err != nil {
			return errors.Join(err, errors.New("config parse function"))
		}
		s.Path = newPath
	}
	f, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(f, t)
	return err
}
