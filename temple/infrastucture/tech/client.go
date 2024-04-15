package tech

import (
	"github.com/ihatiko/olymp/temple/infrastucture/config"
	"github.com/ihatiko/olymp/temple/infrastucture/logger"
	"github.com/ihatiko/olymp/temple/infrastucture/tracer"
)

func Configure(args ...string) error {
	var (
		c   *Config
		err error
	)
	c, err = config.GetConfig[Config]()
	if err != nil {
		return err
	}

	if c.Log != nil {
		c.Log.Configure(
			logger.WithAppName(c.Server.Name),
		)
	}

	if c.Tracer != nil {
		c.Tracer.SetTracer(
			tracer.WithServiceName(c.Server.Name),
			tracer.WithCommand(args[0]),
		)
	}

	return nil
}
