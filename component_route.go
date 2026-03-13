package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
)

type ComponentHandler func(w http.ResponseWriter, r *http.Request) templ.Component

type ComponentRoute struct {
	middlewares []Middleware
	nesters     []Nester
	path        string
	method      string
	handler     ComponentHandler
}

func (c *ComponentRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// create an empty handler
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
