package tech

import (
	_ "embed"

	"github.com/ihatiko/chef/components/observability/http"
	"github.com/ihatiko/chef/components/observability/logger"
	"github.com/ihatiko/chef/components/observability/tracer"
)

//go:embed config/tech.config.toml
var defaultTechConfig []byte

type Service struct {
	Name string `toml:"name"`
}

type Config struct {
	Tech struct {
		Service Service       `toml:"service"`
		Logger  logger.Config `toml:"logger"`
		Tracer  tracer.Config `toml:"tracer"`
		Http    http.Config   `toml:"http"`
	} `toml:"tech"`
}
