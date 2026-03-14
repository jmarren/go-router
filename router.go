package gorouter

type Middleware func(h Handler) Handler

type Router struct {
	// middlewares to apply to all routes
	middlewares []Middleware
	// routes served by the router
	routes []*Route
	// a slice of ErrCatchers that will
	// handle errors returned by route handlers
	catchers []ErrCatcher
}

func CreateRouter() *Router {
	return &Router{
		middlewares: []Middleware{},
		routes:      []*Route{},
	}
}

// adds a middleware to the chain
func (r *Router) Use(m Middleware) {
	r.middlewares = append([]Middleware{m}, r.middlewares...)
}

// adds an ErrCatcher to the routers ErrCatchers array
func (r *Router) UseCatcher(e ErrCatcher) {
	r.catchers = append([]ErrCatcher{e}, r.catchers...)
}

// appends a route with the provided path, method, and handler
func (r *Router) appendRoute(path string, method string, handler Handler) *Route {
	route := &Route{
		path:        path,
		method:      method,
		handler:     handler,
		middlewares: r.middlewares,
		catchers:    r.catchers,
	}
	r.routes = append(r.routes, route)
	return route
}

// registers a get route with the router
func (r *Router) Get(path string, handler Handler) *Route {
	return r.appendRoute(path, "GET", handler)
}

// registers a post route with the router
func (r *Router) Post(path string, handler Handler) *Route {
	return r.appendRoute(path, "POST", handler)
}

// registers a put route with the router
func (r *Router) Put(path string, handler Handler) *Route {
	return r.appendRoute(path, "PUT", handler)
}

// registers a delete route with the router
func (r *Router) Delete(path string, handler Handler) *Route {
	return r.appendRoute(path, "DELETE", handler)
}

/*
Adds a subroute the the router by adding all of its routes.

This includes concatenating the routes path and applying the routers middleware
to the routes handler
*/
func (r *Router) SubRoute(path string, subRoute *Router) {

	for _, route := range subRoute.routes {
		r.routes = append(r.routes, &Route{
			path:        path + route.path,
			method:      route.method,
			handler:     route.handler,
			middlewares: append(route.middlewares, r.middlewares...),
			catchers:    append(route.catchers, r.catchers...),
		})
	}
}
