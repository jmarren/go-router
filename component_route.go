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

	// create a return handler that:
	// - creates component
	// - nests component
	// - renders component
	handler := func(w http.ResponseWriter, r *http.Request) {

		// create the component using the componentHandler
		component := c.componentHandler(w, r)

		// apply nesters to the component
		for i := 0; i < len(c.nesters); i++ {
			component = c.nesters[i](w, r, component)
		}

		// render
		component.Render(r.Context(), w)
	}

	// apply middlewares to the created handler
	// (they will execute before the handler at runtime)
	for i := 0; i < len(c.middlewares); i++ {
		handler = c.middlewares[i](handler)
	}

	return handler
}
