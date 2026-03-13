package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
)

type Middleware func(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)

type Router struct {
	middlewares []Middleware
	routes      []*Route
}

type Nester func(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component

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
