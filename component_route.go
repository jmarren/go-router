package gorouter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/a-h/templ"
	"github.com/jmarren/go-router/views"
)

type ComponentHandler func(rw *RW) (templ.Component, error)

type UnsafeComponentHandler func(rw *RW) templ.Component

type ComponentErrCatcher func(rw *RW, err error) (templ.Component, error)

type Trigger struct {
	event   string
	message string
}

type ComponentRoute struct {
	middlewares          []Middleware
	wrappers             []Wrapper
	path                 string
	method               string
	component            ComponentHandler
	componentErrCatchers []ComponentErrCatcher
	scripts              []string
	triggers             []Trigger
}

func (c *ComponentRoute) Trigger(event, message string) *ComponentRoute {
	fmt.Printf("adding trigger to route = %s: %s\n", event, message)
	c.triggers = append(c.triggers, Trigger{
		event,
		message,
	})
	// c.triggers[event] = message
	return c
}

type IComponentRoute interface {
	Catch(catcher ...ComponentErrCatcher) IComponentRoute
	Use(m Middleware) IComponentRoute
}

func UnsafeComponent(unsafeHandler UnsafeComponentHandler) ComponentHandler {
	return func(rw *RW) (templ.Component, error) {
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

func (c *ComponentRoute) UseScripts(srcs ...string) *ComponentRoute {
	c.scripts = append(c.scripts, srcs...)
	return c
}

func (c *ComponentRoute) head(alreadyExecuted []string) templ.Component {
	toExecute := []string{}
	for _, script := range c.scripts {
		if !slices.Contains(alreadyExecuted, script) {
			toExecute = append(toExecute, script)
		}
	}
	return views.WrapHead(views.ScriptHead(toExecute...))
}

func (c *ComponentRoute) HTTPHandler(baseWrapper baseWrapper) http.HandlerFunc {

	// create a return handler that:
	// - creates component
	// - catches component errors
	// - nests component
	// - renders component
	handler := func(rw *RW) error {

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

		if rw.Request.Header.Get("HX-Request") != "true" {
			component = baseWrapper(component, c.scripts...)
		} else {
			executedStr := rw.Request.Header.Get("HX-Executed")

			var executed []string

			json.Unmarshal([]byte(executedStr), &executed)

			component = templ.Join(component, c.head(executed))
		}

		// create triggers
		if len(c.triggers) > 0 {

			triggerMap := map[string]string{}

			for _, trigger := range c.triggers {
				triggerMap[trigger.event] = trigger.message
			}

			triggersJson, err := json.Marshal(triggerMap)

			if err != nil {
				fmt.Printf("error marshalling triggers into json: %s\n", err)
			} else {
				rw.ResponseWriter.Header().Set("HX-Trigger", string(triggersJson))
			}
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
		err := handler(&RW{
			Request:        r,
			ResponseWriter: w,
		})
		if err != nil {
			panic(err)
		}
	}
}
