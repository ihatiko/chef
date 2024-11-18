package config

import (
	"errors"
	"net/url"
	"os"

	"github.com/ihatiko/olymp/infrastucture/components/utils/toml"
)

const (
	configPath = "config/config.toml"
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
