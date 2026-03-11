package gorouter

import "net/http"

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

type App struct {
	mux *http.ServeMux
	*Router
}

func CreateApp() *App {
	return &App{
		mux:    http.NewServeMux(),
		Router: CreateRouter(),
	}
}

func CreateRouter() *Router {
	return &Router{
		middlewares: []Middleware{},
		routes:      []Route{},
	}
}

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) applyMiddlewares(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (r *Router) Get(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	r.routes = append(r.routes, Route{
		path:    path,
		method:  "GET",
		handler: r.applyMiddlewares(handler),
	})
}
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		path:    path,
		method:  "POST",
		handler: r.applyMiddlewares(handler),
	})
}

func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		path:    path,
		method:  "PUT",
		handler: r.applyMiddlewares(handler),
	})
}

func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		path:    path,
		method:  "DELETE",
		handler: r.applyMiddlewares(handler),
	})
}

func (r *Router) SubRoute(path string, subRoute *Router) {
	for _, route := range subRoute.routes {
		r.routes = append(r.routes, Route{
			path:    path + route.path,
			method:  route.method,
			handler: r.applyMiddlewares(route.handler),
		})
	}
}

//	func (a *App) Get(path string, handler http.HandlerFunc) {
//		a.router.Get(path, handler)
//	}
//
//	func (a *App) Post(path string, handler http.HandlerFunc) {
//		a.router.Post(path, handler)
//	}
//
//	func (a *App) Put(path string, handler http.HandlerFunc) {
//		a.router.Put(path, handler)
//	}
//
//	func (a *App) Delete(path string, handler http.HandlerFunc) {
//		a.router.Delete(path, handler)
//	}
func (a *App) applyRoutes() {
	for _, route := range a.routes {
		a.mux.Handle(route.method+" "+route.path, http.HandlerFunc(route.handler))
	}
}

func (a *App) Serve(addr string) {
	a.applyRoutes()
	http.ListenAndServe(addr, a.mux)
}
