package tech

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
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
	c.Tech.Http.New().Run()
	c.Tech.Log.New()
	c.Tech.Tracer.New()
	return err
}
