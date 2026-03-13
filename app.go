package gorouter

import "net/http"

type App struct {
	mux *http.ServeMux
	*ComponentRouter
}

func CreateApp() *App {
	return &App{
		mux:             http.NewServeMux(),
		ComponentRouter: CreateComponentRouter(),
	}
}

// applies the apps routes to the apps mux
func (a *App) applyRoutes() {

	// handle regular routes with path
	for _, route := range a.routes {
		a.mux.Handle(route.path, route.HTTPHandler())
	}

	// handle component routes with path
	for _, route := range a.componentRoutes {
		a.mux.Handle(route.path, route.HTTPHandler())
	}
}

// applies the apps routes to its mux, then listens on the provided address
func (a *App) Serve(addr string) {
	a.applyRoutes()
	http.ListenAndServe(addr, a.mux)
}
