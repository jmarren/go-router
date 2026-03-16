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
	wrapper Wrapper
	/*
		a slice of all the componentCatchers that will
		catch errors returned by component handlers
	*/
	componentCatchers []ComponentErrCatcher

	scripts []string

	triggers []Trigger

	path string

	prefixWrap bool
}

/* creates an empty ComponentRouter */
func CreateComponentRouter() *ComponentRouter {
	return &ComponentRouter{
		scripts:         []string{},
		Router:          CreateRouter(),
		wrapper:         defaultWrapper(),
		componentRoutes: []*ComponentRoute{},
		triggers:        []Trigger{},
	}
}

func (c *ComponentRouter) Trigger(event, message string) *ComponentRouter {
	c.triggers = append(c.triggers, Trigger{
		event,
		message,
	})
	return c
}

func (c *ComponentRouter) UseScripts(src ...string) *ComponentRouter {
	c.scripts = append(c.scripts, src...)
	return c
}

// func (c *ComponentRouter) UseWrapper(w WrapperFunc) Wrapper {
// 	c.wrapper.UseFunc(w)
// 	return c.wrapper
// }

/*
applies a wrapper to the ComponentRouter so that all subsequently added
componentRoutes are wrapped with the provided function
*/
func (c *ComponentRouter) Wrapper() Wrapper {
	return c.wrapper
}

func (c *ComponentRouter) UsePrefixWrap() *ComponentRouter {
	c.prefixWrap = true
	return c
}

// wraps the component using the provided wrapperFunc only if the
// current url of the request does not contain the provided subpath string
// func (c *ComponentRouter) PrefixWrap(subPath string, w WrapperFunc) Wrapper {
// 	wrapperFunc := func(rw *RW, component templ.Component) (templ.Component, error) {
//
// 		if rw.PathHasPrefix(subPath) {
// 			return component, nil
// 		}
//
// 		return w(rw, component)
// 	}
//
// 	c.wrapper.Use(wrapperFunc)
// 	return c.wrapper
// }

// creates a wrapper with empty err handler,
// applies the hxWrapMiddleware to it,
// then returns it
func (c *ComponentRouter) HxWrap() Wrapper {
	c.wrapper.Use(hxWrapMiddleware)
	return c.wrapper
}

/*
Adds a new componentHandler to the routers routes with the provided path and method

The newly added route inherits the properties of the router (middlewares, catchers, wrappers)

A pointer to the added route is returned so that methods may be chained
*/
func (c *ComponentRouter) addComponentRoute(path string, ch ComponentHandler, method string) *ComponentRoute {

	route := &ComponentRoute{
		wrapper:              c.wrapper,
		path:                 path,
		method:               method,
		component:            ch,
		middlewares:          c.middlewares,
		componentErrCatchers: c.componentCatchers,
		scripts:              c.scripts,
		triggers:             c.triggers,
		shouldWrap:           true,
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

	subComponent.wrapper.UseFunc(c.wrapper.wrapperFunc())
	if subComponent.prefixWrap {
		subComponent.wrapper.Use(PrefixWrap(path))
	}

	for _, cr := range subComponent.componentRoutes {

		// wrapper :=

		newRoute := &ComponentRoute{
			path:                 path + cr.path,
			method:               cr.method,
			component:            cr.component,
			wrapper:              cr.wrapper,
			middlewares:          append(cr.middlewares, c.middlewares...),
			componentErrCatchers: append(cr.componentErrCatchers, c.componentCatchers...),
			scripts:              append(cr.scripts, c.scripts...),
			triggers:             append(cr.triggers, c.triggers...),
			shouldWrap:           cr.shouldWrap,
		}
		// copy from c to cr so that component triggers overwrite router triggers on conflict
		c.componentRoutes = append(c.componentRoutes, newRoute)

	}

	// add the subComponents regular router as a subroute as well
	c.SubRoute(path, subComponent.Router)
}
