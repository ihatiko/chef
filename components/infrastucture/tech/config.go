package tech

import (
	_ "embed"

	"github.com/ihatiko/olymp/components/clients/http"
	"github.com/ihatiko/olymp/components/clients/logger"
	"github.com/ihatiko/olymp/components/clients/tracer"
)

//go:embed config/tech.config.toml
var defaultConfig []byte

type Service struct {
	Name string `toml:"name"`
}

type Config struct {
	Tech struct {
		Service Service       `toml:"service"`
		Log     logger.Config `toml:"log"`
		Tracer  tracer.Config `toml:"tracer"`
		Http    http.Config   `toml:"http"`
	} `toml:"tech"`
}
