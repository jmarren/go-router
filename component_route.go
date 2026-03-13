package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
)

type ComponentHandler func(w http.ResponseWriter, r *http.Request) templ.Component

type ComponentRoute struct {
	middlewares      []Middleware
	nesters          []Nester
	path             string
	method           string
	componentHandler ComponentHandler
}

func (c *ComponentRoute) HTTPHandler() http.HandlerFunc {

	// create an empty default handler
	handler := func(w http.ResponseWriter, r *http.Request) {}

	// apply middlewares to handler
	for i := 0; i < len(c.middlewares); i++ {
		handler = c.middlewares[i](handler)
	}

	// create a return handler that:
	// - executes middlewares
	// - creates component
	// - nests component
	// - renders component
	ret := func(w http.ResponseWriter, r *http.Request) {
		// apply the handler so that middlewares are executed
		// handler(w, r)
		handler(w, r)

		// create the component using the componentHandler
		component := c.componentHandler(w, r)

		// apply nesters to the component
		for i := 0; i < len(c.nesters); i++ {
			component = c.nesters[i](w, r, component)
		}

		// render
		component.Render(r.Context(), w)
	}

	return http.HandlerFunc(ret)

}

// func (c *ComponentRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	// create an empty handler
// 	handler := func(w http.ResponseWriter, r *http.Request) {}
//
// 	// apply middlewares to handler
// 	for i := 0; i < len(c.middlewares); i++ {
// 		handler = c.middlewares[i](handler)
// 	}
//
// 	// apply the handler so that middlewares are executed
// 	handler(w, r)
//
// 	component := c.componentHandler(w, r)
//
// 	// apply nesters
// 	for i := 0; i < len(c.nesters); i++ {
// 		component = c.nesters[i](w, r, component)
// 	}
//
// 	component.Render(r.Context(), w)
// }
