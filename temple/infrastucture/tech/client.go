package tech

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func Use(arg string) error {
	c := new(Config)
	err := toml.Unmarshal(defaultConfig, c)
	if err != nil {
		e := fmt.Errorf("Error unmarshalling tech-config: %s command", err, arg)
		fmt.Println("-----------------------")
		fmt.Println(e)
		fmt.Println("-----------------------")
		return e
	}
	//TODO env rewrite
	c.Tech.Http.Use().Run()
	c.Tech.Log.Use()
	c.Tech.Tracer.Use()
	return err
}
