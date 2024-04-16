package tech

import (
	"github.com/ihatiko/olymp/temple/infrastucture/http"
	"github.com/ihatiko/olymp/temple/infrastucture/logger"
	"github.com/ihatiko/olymp/temple/infrastucture/tracer"
)

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
