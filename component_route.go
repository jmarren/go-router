package gorouter

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

type ComponentHandler func(w http.ResponseWriter, r *http.Request) (templ.Component, error)

type ComponentErrCatcher func(w http.ResponseWriter, r *http.Request, err error) (templ.Component, error)

type ComponentRoute struct {
	middlewares []Middleware
	wrappers    []ComponentWrapper
	path        string
	method      string
	component   ComponentHandler
	errCatchers []ComponentErrCatcher
}

type IComponentRoute interface {
	Catch(catcher ComponentErrCatcher) IComponentRoute
}

func UnsafeComponent(unsafeHandler func(w http.ResponseWriter, r *http.Request) templ.Component) ComponentHandler {
	return func(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
		return unsafeHandler(w, r), nil
	}
}

func (c *ComponentRoute) Catch(catcher ComponentErrCatcher) IComponentRoute {
	c.errCatchers = append([]ComponentErrCatcher{catcher}, c.errCatchers...)
	return c
}

func (c *ComponentRoute) HTTPHandler() http.HandlerFunc {

	// create a return handler that:
	// - creates component
	// - nests component
	// - renders component
	handler := func(w http.ResponseWriter, r *http.Request) {

		// create the component using the componentHandler
		component, err := c.component(w, r)

		fmt.Printf("num errCatchers = %d\n", len(c.errCatchers))

		if err != nil {
			for _, catcher := range c.errCatchers {
				component, err = catcher(w, r, err)
				if err == nil {
					break
				}
			}
		}

		// panic if the error is still not resolved
		if err != nil {
			panic(err)
		}

		// apply nesters to the component
		for i := 0; i < len(c.wrappers); i++ {
			component, err = c.wrappers[i].wrap(w, r, component)
			if err != nil {
				component, err = c.wrappers[i].err(w, r, err)
			}
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
