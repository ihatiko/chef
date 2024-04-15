package tech

import (
	"github.com/ihatiko/config"
	"github.com/ihatiko/olymp/logger"
	"github.com/ihatiko/olymp/tracer"
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
