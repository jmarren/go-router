package gorouter

type ComponentRouter struct {
	*Router
	componentRoutes []*ComponentRoute
	wrappers        []ComponentWrapper
	errCatchers     []ComponentErrCatcher
}

func CreateComponentRouter() *ComponentRouter {
	return &ComponentRouter{
		Router:      CreateRouter(),
		wrappers:    []ComponentWrapper{},
		errCatchers: []ComponentErrCatcher{},
	}
}

func (c *ComponentRouter) UseNester(n ComponentWrapper) {
	c.wrappers = append([]ComponentWrapper{n}, c.wrappers...)
}

func (c *ComponentRouter) UseHxNester(n SimpleNester) {
	c.UseNester(UnsafeHxReqWrapper(n))
}

func (c *ComponentRouter) UseComponentCatcher(ec ComponentErrCatcher) {
	c.errCatchers = append([]ComponentErrCatcher{ec}, c.errCatchers...)
}

func (c *ComponentRouter) appendComponentRoute(path string, ch ComponentHandler, method string) *ComponentRoute {
	route := &ComponentRoute{
		wrappers:    c.wrappers,
		path:        path,
		method:      method,
		component:   ch,
		middlewares: c.middlewares,
		errCatchers: c.errCatchers,
	}

	c.componentRoutes = append(c.componentRoutes, route)
	return route
}

func (c *ComponentRouter) GetComponent(path string, ch ComponentHandler) *ComponentRoute {
	return c.appendComponentRoute(path, ch, "GET")
}

func (c *ComponentRouter) PostComponent(path string, ch ComponentHandler) *ComponentRoute {
	return c.appendComponentRoute(path, ch, "POST")
}

func (c *ComponentRouter) PutComponent(path string, ch ComponentHandler) *ComponentRoute {
	return c.appendComponentRoute(path, ch, "PUT")
}

func (c *ComponentRouter) DeleteComponent(path string, ch ComponentHandler) *ComponentRoute {
	return c.appendComponentRoute(path, ch, "DELETE")
}

func (c *ComponentRouter) SubComponent(path string, subComponent *ComponentRouter) {
	for _, cr := range subComponent.componentRoutes {
		c.componentRoutes = append(c.componentRoutes, &ComponentRoute{
			path:        path + cr.path,
			method:      cr.method,
			component:   cr.component,
			wrappers:    append(cr.wrappers, c.wrappers...),
			middlewares: append(cr.middlewares, c.middlewares...),
			errCatchers: append(cr.errCatchers, c.errCatchers...),
		})
	}

	// add the subComponents regular router as a subroute as well
	c.SubRoute(path, subComponent.Router)
}
