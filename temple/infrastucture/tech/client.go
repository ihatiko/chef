package tech

func Use() {
	var (
		c Config
	)

	c.Tech.Http.Use()
	c.Tech.Log.Use()
	c.Tech.Tracer.Use()
}
