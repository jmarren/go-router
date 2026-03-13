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

func (c *ComponentRouter) GetComponent(path string, ch ComponentHandler) {
	c.componentRoutes = append(c.componentRoutes, &ComponentRoute{
		nesters:     c.nesters,
		path:        path,
		method:      "GET",
		handler:     ch,
		middlewares: c.middlewares,
	})
}

func (c *ComponentRouter) PostComponent(path string, ch ComponentHandler) {
	c.componentRoutes = append(c.componentRoutes, &ComponentRoute{
		nesters:     c.nesters,
		path:        path,
		method:      "POST",
		handler:     ch,
		middlewares: c.middlewares,
	})
}

func (c *ComponentRouter) PutComponent(path string, ch ComponentHandler) {
	c.componentRoutes = append(c.componentRoutes, &ComponentRoute{
		nesters:     c.nesters,
		path:        path,
		method:      "PUT",
		handler:     ch,
		middlewares: c.middlewares,
	})
}

func (c *ComponentRouter) DeleteComponent(path string, ch ComponentHandler) {
	c.componentRoutes = append(c.componentRoutes, &ComponentRoute{
		path:        path,
		nesters:     c.nesters,
		method:      "DELETE",
		handler:     ch,
		middlewares: c.middlewares,
	})
}

func (c *ComponentRouter) SubComponent(path string, subComponent *ComponentRouter) {
	for _, cr := range subComponent.componentRoutes {
		c.componentRoutes = append(c.componentRoutes, &ComponentRoute{
			path:        path + cr.path,
			method:      cr.method,
			handler:     cr.handler,
			nesters:     append(cr.nesters, c.nesters...),
			middlewares: append(cr.middlewares, c.middlewares...),
		})
	}

	// add the subComponents regular router as a subroute as well
	c.SubRoute(path, subComponent.Router)
}
