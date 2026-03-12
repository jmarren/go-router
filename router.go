package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/jmarren/go-router/views"
)

type Middleware func(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request)

type Route struct {
	path    string
	method  string
	handler func(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	middlewares []Middleware
	routes      []Route
}

type Nester func(w http.ResponseWriter, r *http.Request) bool

type ComponentRouter struct {
	*Router
	base   BaseComponent
	nester Nester
}

type BaseComponent func(base templ.Component) templ.Component

type App struct {
	mux *http.ServeMux
	// *Router
	*ComponentRouter
	// base BaseComponent
}

func (c *ComponentRouter) UseBase(base BaseComponent) {
	c.base = base
}

func (c *ComponentRouter) UseNester(n Nester) {
	c.nester = n
}

type ComponentHandler func(w http.ResponseWriter, r *http.Request) templ.Component

func (c *ComponentRouter) nest(ch ComponentHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		component := ch(w, r)

		// if nester returns true
		// wrap the component with base
		if c.nester(w, r) {
			component = c.base(component)
		}
		component.Render(r.Context(), w)
	}

}
func (c *ComponentRouter) GetComponent(path string, ch ComponentHandler) {
	c.Get(path, c.nest(ch))
}

func (c *ComponentRouter) PostComponent(path string, ch ComponentHandler) {
	c.Post(path, c.nest(ch))
}

func (c *ComponentRouter) PutComponent(path string, ch ComponentHandler) {
	c.Put(path, c.nest(ch))
}

func (c *ComponentRouter) DeleteComponent(path string, ch ComponentHandler) {
	c.Delete(path, c.nest(ch))
}

// func (a *App)
//
// type ShouldRenderFull func(w http.ResponseWriter, r *http.Request) bool
//
// type NestedComponentRouter struct {
// 	base             BaseComponent
// 	shouldRenderFull ShouldRenderFull
// 	*Router
// }
//
// func CreateNestedComponentRouter(baseFn func(base templ.Component) templ.Component, shouldRenderFull ShouldRenderFull) *NestedComponentRouter {
// 	return &NestedComponentRouter{
// 		base:             baseFn,
// 		Router:           CreateRouter(),
// 		shouldRenderFull: shouldRenderFull,
// 	}
// }
//
// func (n *NestedComponentRouter) handler(ch func(w http.ResponseWriter, r *http.Request) templ.Component) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		component := ch(w, r)
// 		if n.shouldRenderFull(w, r) {
// 			component = n.base(component)
// 		}
// 		component.Render(r.Context(), w)
// 	}
// }
//
// func (n *NestedComponentRouter) Get(path string, ch func(w http.ResponseWriter, r *http.Request) templ.Component) {
// 	n.Router.Get(path, n.handler(ch))
// }

//	type ComponentRouter struct {
//		Get    func()
//		Post   func()
//		Put    func()
//		Delete func()
//	}

func CreateRouter() *Router {
	return &Router{
		middlewares: []Middleware{},
		routes:      []Route{},
	}
}

func CreateComponentRouter() *ComponentRouter {
	return &ComponentRouter{
		Router: CreateRouter(),
		base: func(base templ.Component) templ.Component {
			return base
		},
		nester: func(w http.ResponseWriter, r *http.Request) bool {
			return true
		},
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
	r.middlewares = append(r.middlewares, m)
}

// applies a routers middleware to the provided handler
func (r *Router) applyMiddlewares(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	return handler
}

// registers a get route with the router
func (r *Router) Get(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.routes = append(r.routes, Route{
		path:    path,
		method:  "GET",
		handler: r.applyMiddlewares(handler),
	})
}

// registers a post route with the router
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		path:    path,
		method:  "POST",
		handler: r.applyMiddlewares(handler),
	})
}

// registers a put route with the router
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		path:    path,
		method:  "PUT",
		handler: r.applyMiddlewares(handler),
	})
}

// registers a delete route with the router
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		path:    path,
		method:  "DELETE",
		handler: r.applyMiddlewares(handler),
	})
}

// renders a component for the given get request path
// func (r *Router) Component(path string, ch func(w http.ResponseWriter, r *http.Request) templ.Component) *ComponentRouter {
//
//		handler := func(w http.ResponseWriter, r *http.Request) {
//			component := ch(w, r)
//			component.Render(r.Context(), w)
//		}
//		return &ComponentRouter{
//			Get: func() {
//				r.Get(path, handler)
//			},
//			Post: func() {
//				r.Post(path, handler)
//			},
//			Put: func() {
//				r.Put(path, handler)
//			},
//			Delete: func() {
//				r.Delete(path, handler)
//			},
//		}
//	}
func (r *Router) Page(path string, method string, ch func(w http.ResponseWriter, r *http.Request) templ.Component) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		component := ch(w, r)
		// if hx-request render only the component
		// otherwise render inside the root page component
		if r.Header.Get("HX-Request") == "true" {
			component.Render(r.Context(), w)
		} else {
			views.Page(component).Render(r.Context(), w)
		}
	}
	r.Get(path, handler)
}

/*
Adds a subroute the the router by adding all of its routes.

This includes concatenating the routes path and applying the routers middleware
to the routes handler
*/
func (r *Router) SubRoute(path string, subRoute *Router) {
	for _, route := range subRoute.routes {
		r.routes = append(r.routes, Route{
			path:    path + route.path,
			method:  route.method,
			handler: r.applyMiddlewares(route.handler),
		})
	}
}

// applies the apps routes to the apps mux
func (a *App) applyRoutes() {
	for _, route := range a.routes {
		a.mux.Handle(route.method+" "+route.path, http.HandlerFunc(route.handler))
	}
}

// applies the apps routes to its mux, then listens on the provided address
func (a *App) Serve(addr string) {
	a.applyRoutes()
	http.ListenAndServe(addr, a.mux)
}
