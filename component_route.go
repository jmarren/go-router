package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
)

type ComponentHandler func(w http.ResponseWriter, r *http.Request) (templ.Component, error)

type UnsafeComponentHandler func(w http.ResponseWriter, r *http.Request) templ.Component

type ComponentErrCatcher func(w http.ResponseWriter, r *http.Request, err error) (templ.Component, error)

type ComponentRoute struct {
	middlewares        []Middleware
	wrapperMiddlewares []WrapMiddleware
	// wrappers             []ComponentWrapper
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
	return func(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
		return unsafeHandler(w, r), nil
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

	wrapper := func(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error) {
		return component, nil
	}

	for _, wm := range c.wrapperMiddlewares {
		wrapper = wm(wrapper)
	}

	// create a return handler that:
	// - creates component
	// - catches component errors
	// - nests component
	// - renders component
	handler := func(w http.ResponseWriter, r *http.Request) error {

		// create the component using the componentHandler
		component, err := c.component(w, r)

		// if an error occurs,
		// apply catchers until it is resolved to nil
		if err != nil {
			for _, catcher := range c.componentErrCatchers {
				component, err = catcher(w, r, err)
				if err == nil {
					break
				}
			}
		}

		if err != nil {
			return err
		}

		// wrap the component
		component, err = wrapper(w, r, component)

		if err != nil {
			return err
		}

		// render
		component.Render(r.Context(), w)

		return nil
	}

	// apply middlewares to the created handler
	// (they will execute before the handler at runtime)
	for i := 0; i < len(c.middlewares); i++ {
		handler = c.middlewares[i](handler)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			panic(err)
		}
	}
}
