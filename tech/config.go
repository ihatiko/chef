package tech

import (
	"github.com/ihatiko/olymp/logger"
	"github.com/ihatiko/olymp/tracer"
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
