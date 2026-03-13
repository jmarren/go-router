package gorouter

import (
	"net/http"
)

type Middleware func(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)

type Router struct {
	middlewares []Middleware
	routes      []*Route
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

// appends a route with the provided path, method, and handler
func (r *Router) appendRoute(path string, method string, handler http.HandlerFunc) {
	r.routes = append(r.routes, &Route{
		path:        path,
		method:      method,
		handler:     handler,
		middlewares: r.middlewares,
	})
}

// registers a get route with the router
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.appendRoute(path, "GET", handler)
}

// registers a post route with the router
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.appendRoute(path, "POST", handler)
}

// registers a put route with the router
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.appendRoute(path, "PUT", handler)
}

// registers a delete route with the router
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.appendRoute(path, "DELETE", handler)
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
