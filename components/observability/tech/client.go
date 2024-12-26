package tech

import (
	"fmt"
	tC "github.com/ihatiko/chef/components/tech/config"
	"github.com/ihatiko/chef/components/tech/toml"
)

func Use(arg string) error {
	c := new(Config)
	err := toml.Unmarshal(defaultTechConfig, c)
	if err != nil {
		e := fmt.Errorf("error unmarshalling tech-config: %v command %s", err, arg)
		return e
	}
	err = tC.ToConfig(c)
	if err != nil {
		e := fmt.Errorf("error unmarshalling tech-config (applyed main config): %v command %s", err, arg)
		return e
	}
	c.Tech.Http.New().Run()
	c.Tech.Logger.New()
	c.Tech.Tracer.New()
	return err
}
