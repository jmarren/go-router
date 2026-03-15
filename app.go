package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
)

type baseWrapper func(component templ.Component, scripts ...string) templ.Component

type App struct {
	mux *http.ServeMux
	*ComponentRouter
	baseWrapper baseWrapper
}

func CreateApp() *App {
	return &App{
		mux:             http.NewServeMux(),
		ComponentRouter: CreateComponentRouter(),
	}
}

func (a *App) UseBaseWrapper(bw baseWrapper) {
	a.baseWrapper = bw
}

func (a *App) UseStaticDir(dir string) {

	// Create a file server handler
	handler := http.StripPrefix("/static/", http.FileServer(http.Dir(dir)))

	// Handle requests at the root URL ("/") using the file server
	// http.Handle("/", handler)
	a.mux.Handle("/static/", handler)
}

// applies the apps routes to the apps mux
func (a *App) applyRoutes() {

	// handle regular routes with path
	for _, route := range a.routes {
		a.mux.Handle(route.path, route.HTTPHandler())
	}

	// handle component routes with path
	for _, route := range a.componentRoutes {
		a.mux.Handle(route.path, route.HTTPHandler(a.baseWrapper))
	}
}

// applies the apps routes to its mux, then listens on the provided address
func (a *App) Serve(addr string) {
	a.applyRoutes()
	http.ListenAndServe(addr, a.mux)
}
