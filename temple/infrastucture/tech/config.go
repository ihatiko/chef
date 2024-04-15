package tech

import (
	"github.com/ihatiko/olymp/temple/infrastucture/logger"
	"github.com/ihatiko/olymp/temple/infrastucture/tracer"
)

type Server struct {
	Name string
	Team string
}
type Config struct {
	Server *Server
	Log    *logger.Config
	Tracer *tracer.Config
}
