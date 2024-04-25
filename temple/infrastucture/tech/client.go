package tech

func Use() {
	var (
		c Config
	)
	//TODO load tech config
	c.Tech.Http.Use()
	c.Tech.Log.Use()
	c.Tech.Tracer.Use()
}
