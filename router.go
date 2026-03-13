package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
)

type Middleware func(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)

type Route struct {
	path        string
	method      string
	middlewares []Middleware
	handler     func(w http.ResponseWriter, r *http.Request)
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(route.middlewares); i++ {
		route.handler = route.middlewares[i](route.handler)
	}
	route.handler(w, r)
}

type ComponentRoute struct {
	middlewares []Middleware
	nesters     []Nester
	path        string
	method      string
	handler     func(w http.ResponseWriter, r *http.Request) templ.Component
}

func (c *ComponentRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	// apply middlewares
	for i := 0; i < len(c.middlewares); i++ {
		handler = c.middlewares[i](handler)
	}

	component := c.handler(w, r)

	// apply nesters
	for i := 0; i < len(c.nesters); i++ {
		component = c.nesters[i](w, r, component)
	}

	component.Render(r.Context(), w)
}

type Router struct {
	middlewares []Middleware
	routes      []*Route
}

type Nester func(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component

type ComponentRouter struct {
	*Router
	componentRoutes []*ComponentRoute
	nesters         []Nester
}

func (c *ComponentRouter) UseNester(n Nester) {
	c.nesters = append([]Nester{n}, c.nesters...)
}

type ComponentHandler func(w http.ResponseWriter, r *http.Request) templ.Component

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
func CreateRouter() *Router {
	return &Router{
		middlewares: []Middleware{},
		routes:      []*Route{},
	}
}

func CreateComponentRouter() *ComponentRouter {
	return &ComponentRouter{
		Router:  CreateRouter(),
		nesters: []Nester{},
	}
}

func CreateApp() *App {
	return &App{
		mux:             http.NewServeMux(),
		ComponentRouter: CreateComponentRouter(),
	}
}

// adds a middleware to the chain
func (r *Router) Use(m Middleware) {
	r.middlewares = append([]Middleware{m}, r.middlewares...)
}

// registers a get route with the router
func (r *Router) Get(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.routes = append(r.routes, &Route{
		path:        path,
		method:      "GET",
		handler:     handler,
		middlewares: r.middlewares,
	})
}

// registers a post route with the router
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, &Route{
		path:        path,
		method:      "POST",
		handler:     handler,
		middlewares: r.middlewares,
	})
}

// registers a put route with the router
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, &Route{
		path:        path,
		method:      "PUT",
		handler:     handler,
		middlewares: r.middlewares,
	})
}

// registers a delete route with the router
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, &Route{
		path:        path,
		method:      "DELETE",
		handler:     handler,
		middlewares: r.middlewares,
	})
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
		})
	}
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

	// add the subComponents router as a subroute as well
	c.SubRoute(path, subComponent.Router)
}
