package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
)

type ComponentHandler func(rw RW) (templ.Component, error)

type UnsafeComponentHandler func(rw RW) templ.Component

type ComponentErrCatcher func(rw RW, err error) (templ.Component, error)

type ComponentRoute struct {
	middlewares          []Middleware
	wrappers             []Wrapper
	path                 string
	method               string
	component            ComponentHandler
	componentErrCatchers []ComponentErrCatcher
}

type IComponentRoute interface {
	Catch(catcher ...ComponentErrCatcher) IComponentRoute
	Use(m Middleware) IComponentRoute
}

func UnsafeComponent(unsafeHandler UnsafeComponentHandler) ComponentHandler {
	return func(rw RW) (templ.Component, error) {
		return unsafeHandler(rw), nil
	}
}

func (c *ComponentRoute) Catch(catcher ...ComponentErrCatcher) IComponentRoute {
	c.componentErrCatchers = append(catcher, c.componentErrCatchers...)
	return c
}

func (c *ComponentRoute) Use(m Middleware) IComponentRoute {
	c.middlewares = append([]Middleware{m}, c.middlewares...)
	return c
}

func (c *ComponentRoute) HTTPHandler() http.HandlerFunc {

	// create a return handler that:
	// - creates component
	// - catches component errors
	// - nests component
	// - renders component
	handler := func(rw RW) error {
		// rw := RW{
		// 	ResponseWriter: w,
		// 	Request:        r,
		// }

		// create the component using the componentHandler
		component, err := c.component(rw)

		// if an error occurs,
		// apply catchers until it is resolved to nil
		if err != nil {
			for _, catcher := range c.componentErrCatchers {
				component, err = catcher(rw, err)
				if err == nil {
					break
				}
			}
		}

		if err != nil {
			return err
		}

		// wrap the component
		for _, wrapper := range c.wrappers {
			// attempt to wrap
			component, err = wrapper.wrap(rw, component)
			// if an err is returned attempt to resolve with err method
			if err != nil {
				component, err = wrapper.err(rw, component, err)
			}
			// if error is unresolved, return it
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

		// render
		component.Render(rw.Request.Context(), rw.ResponseWriter)

		return nil
	}

	// apply middlewares to the created handler
	// (they will execute before the handler at runtime)
	for i := 0; i < len(c.middlewares); i++ {
		handler = c.middlewares[i](handler)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(RW{
			Request:        r,
			ResponseWriter: w,
		})
		if err != nil {
			panic(err)
		}
	}
}
