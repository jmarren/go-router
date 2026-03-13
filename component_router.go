package gorouter

type ComponentRouter struct {
	*Router
	componentRoutes []*ComponentRoute
	nesters         []Nester
}

func CreateComponentRouter() *ComponentRouter {
	return &ComponentRouter{
		Router:  CreateRouter(),
		nesters: []Nester{},
	}
}

func (c *ComponentRouter) UseNester(n Nester) {
	c.nesters = append([]Nester{n}, c.nesters...)
}

func (c *ComponentRouter) UseHxNester(n SimpleNester) {
	c.UseNester(HxReqNester(n))
}

func (c *ComponentRouter) appendComponentRoute(path string, ch ComponentHandler, method string) {
	c.componentRoutes = append(c.componentRoutes, &ComponentRoute{
		nesters:          c.nesters,
		path:             path,
		method:           method,
		componentHandler: ch,
		middlewares:      c.middlewares,
	})
}

func (c *ComponentRouter) GetComponent(path string, ch ComponentHandler) {
	c.appendComponentRoute(path, ch, "GET")
}

func (c *ComponentRouter) PostComponent(path string, ch ComponentHandler) {
	c.appendComponentRoute(path, ch, "POST")
}

func (c *ComponentRouter) PutComponent(path string, ch ComponentHandler) {
	c.appendComponentRoute(path, ch, "PUT")
}

func (c *ComponentRouter) DeleteComponent(path string, ch ComponentHandler) {
	c.appendComponentRoute(path, ch, "DELETE")
}

func (c *ComponentRouter) SubComponent(path string, subComponent *ComponentRouter) {
	for _, cr := range subComponent.componentRoutes {
		c.componentRoutes = append(c.componentRoutes, &ComponentRoute{
			path:             path + cr.path,
			method:           cr.method,
			componentHandler: cr.componentHandler,
			nesters:          append(cr.nesters, c.nesters...),
			middlewares:      append(cr.middlewares, c.middlewares...),
		})
	}

	// add the subComponents regular router as a subroute as well
	c.SubRoute(path, subComponent.Router)
}
