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
	Name string
}

type Config struct {
	Tech struct {
		Service Service
		Log     logger.Config
		Tracer  tracer.Config
		Http    http.Config
	}
}
