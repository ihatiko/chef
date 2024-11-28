package tech

import (
	"fmt"
	"github.com/ihatiko/olymp/infrastucture/components/utils/toml"
	"go.uber.org/zap"
	"log/slog"
)

func Use(arg string) error {
	c := new(Config)
	err := toml.Unmarshal(defaultConfig, c)
	if err != nil {
		e := fmt.Errorf("error unmarshalling tech-config: %s command %s", err, arg)
		slog.Error("error unmarshalling tech-config", slog.String("error", e.Error()), zap.String("command", arg))
		return e
	}
	//TODO env rewrite
	c.Tech.Http.New().Run()
	c.Tech.Log.New()
	c.Tech.Tracer.New()
	return err
}
