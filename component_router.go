package gorouter

/* A router that serves components */
type ComponentRouter struct {
	/*
		embeds the Router struct so that regular routes may be used
		alongside component routes
	*/
	*Router
	/*
		a slice of all the component routes to be served by this router
	*/
	componentRoutes []*ComponentRoute
	/*
		a slice of all the wrappers that will wrap components served
		by the router
	*/
	wrappersMiddlewares []WrapMiddleware
	/*
		a slice of all the componentCatchers that will
		catch errors returned by component handlers
	*/
	componentCatchers []ComponentErrCatcher
}

/* creates an empty ComponentRouter */
func CreateComponentRouter() *ComponentRouter {
	return &ComponentRouter{
		Router:              CreateRouter(),
		wrappersMiddlewares: []WrapMiddleware{},
		componentRoutes:     []*ComponentRoute{},
	}
}

/*
applies a wrapper to the ComponentRouter so that all subsequently added
componentRoutes are wrapped with the provided function
*/
func (c *ComponentRouter) UseWrapperMiddleware(w WrapMiddleware) {
	c.wrappersMiddlewares = append([]WrapMiddleware{w}, c.wrappersMiddlewares...)
}

/*
applies the unsafeHxWrapper so that wrapping occurs if the request has
the HX-Request header.

Does not handle errors
*/
func (c *ComponentRouter) SimpleWrapper(n SimpleWrapper) {
	c.UseWrapperMiddleware(MiddlewareFromSimple(n))
}

func (c *ComponentRouter) UseSimpleHxWrapper(n SimpleWrapper) {
	c.UseHxWrapper(FromSimple(n))
}

func (c *ComponentRouter) UseHxWrapper(w Wrapper) {
	c.UseWrapperMiddleware(hxWrapMiddleware)
	c.UseWrapper(w)
}

func (c *ComponentRouter) UseWrapper(wrapper Wrapper) {
	// convert the wrapper to a middleware
	wrapMiddleware := func(w Wrapper) Wrapper {
		return wrapper
	}
	c.wrappersMiddlewares = append([]WrapMiddleware{wrapMiddleware}, c.wrappersMiddlewares...)
}

/*
Adds a new componentHandler to the routers routes with the provided path and method

The newly added route inherits the properties of the router (middlewares, catchers, wrappers)

A pointer to the added route is returned so that methods may be chained
*/
func (c *ComponentRouter) addComponentRoute(path string, ch ComponentHandler, method string) *ComponentRoute {
	route := &ComponentRoute{
		wrapperMiddlewares:   c.wrappersMiddlewares,
		path:                 path,
		method:               method,
		component:            ch,
		middlewares:          c.middlewares,
		componentErrCatchers: c.componentCatchers,
	}

	c.componentRoutes = append(c.componentRoutes, route)
	return route
}

/* HTTP METHODS */
func (c *ComponentRouter) GetComponent(path string, ch ComponentHandler) *ComponentRoute {
	return c.addComponentRoute(path, ch, "GET")
}

func (c *ComponentRouter) PostComponent(path string, ch ComponentHandler) *ComponentRoute {
	return c.addComponentRoute(path, ch, "POST")
}

func (c *ComponentRouter) PutComponent(path string, ch ComponentHandler) *ComponentRoute {
	return c.addComponentRoute(path, ch, "PUT")
}

func (c *ComponentRouter) DeleteComponent(path string, ch ComponentHandler) *ComponentRoute {
	return c.addComponentRoute(path, ch, "DELETE")
}

/*
The SubComponent method mounts another component router onto this one.

The mounted component inherits all the properties of the mounter
*/
func (c *ComponentRouter) SubComponent(path string, subComponent *ComponentRouter) {
	for _, cr := range subComponent.componentRoutes {
		c.componentRoutes = append(c.componentRoutes, &ComponentRoute{
			path:               path + cr.path,
			method:             cr.method,
			component:          cr.component,
			wrapperMiddlewares: append(cr.wrapperMiddlewares, c.wrappersMiddlewares...),
			middlewares:        append(cr.middlewares, c.middlewares...),
			// componentErrCatchers: append(cr.componentErrCatchers, c.componentCatchers...),
		})
	}

	// add the subComponents regular router as a subroute as well
	c.SubRoute(path, subComponent.Router)
}
