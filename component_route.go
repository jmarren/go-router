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
	wrapper              Wrapper
	path                 string
	method               string
	component            ComponentHandler
	componentErrCatchers []ComponentErrCatcher
	scripts              []string
	triggers             []Trigger
	shouldWrap           bool
	target               string
}

func (c *ComponentRoute) Trigger(event, message string) *ComponentRoute {
	c.triggers = append(c.triggers, Trigger{
		event,
		message,
	})
	return c
}

func SimpleComponent(componentFunc func() templ.Component) ComponentHandler {
	return func(rw *RW) (templ.Component, error) {
		return componentFunc(), nil
	}
}

func UnsafeComponent(unsafeHandler UnsafeComponentHandler) ComponentHandler {
	return func(rw *RW) (templ.Component, error) {
		return unsafeHandler(rw), nil
	}
}

func (c *ComponentRoute) Catch(catcher ...ComponentErrCatcher) *ComponentRoute {
	c.componentErrCatchers = append(catcher, c.componentErrCatchers...)
	return c
}

func (c *ComponentRoute) Retarget(target string) *ComponentRoute {
	c.target = target
	return c
}

func (c *ComponentRoute) DontWrap() *ComponentRoute {
	c.shouldWrap = false
	return c
}

func (c *ComponentRoute) Use(m Middleware) *ComponentRoute {
	c.middlewares = append([]Middleware{m}, c.middlewares...)
	return c
}

func (c *ComponentRoute) UseScripts(srcs ...string) *ComponentRoute {
	c.scripts = append(c.scripts, srcs...)
	return c
}

func (c *ComponentRoute) head(alreadyExecuted []string) templ.Component {
	toExecute := []string{}

	// only add unexecuted scripts to the toExecute slice
	for _, script := range c.scripts {
		if !slices.Contains(alreadyExecuted, script) {
			toExecute = append(toExecute, script)
		}
	}
	return views.WrapHead(views.ScriptHead(toExecute...))
}

func (c *ComponentRoute) triggersJson() string {

	if len(c.triggers) == 0 {
		return ""
	}

	// create a map to marshal into json properly
	triggerMap := map[string]string{}

	for _, trigger := range c.triggers {
		triggerMap[trigger.event] = trigger.message
	}

	triggersJson, err := json.Marshal(triggerMap)

	if err != nil {
		fmt.Printf("error marshalling triggers into json: %s\n", err)
	}

	return string(triggersJson)

}

func (c *ComponentRoute) HTTPHandler(baseWrapper baseWrapper) http.HandlerFunc {

	// create a return handler that:
	// - creates component
	// - catches component errors
	// - nests component
	// - renders component
	handler := func(rw *RW) error {

		rw.Retarget(c.target)

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

			if err != nil {
				return err
			}
		}

		if c.shouldWrap {
			// wrap the component
			component, err = c.wrapper.wrapperFunc()(rw, component)
			if err != nil {
				return err
			}
		}

		// add scripts
		if !rw.IsHxRequest() {
			component = baseWrapper(component, c.scripts...)
		} else {
			executed := rw.ExecutedScripts()
			component = templ.Join(component, c.head(executed))
		}

		// // retarget
		if rw.target != "" {
			fmt.Printf("rw.target =  %s\n", rw.target)
			rw.ResponseWriter.Header().Set("HX-Retarget", rw.target)
		}

		// add triggers
		triggersJson := c.triggersJson()

		if triggersJson != "" {
			rw.ResponseWriter.Header().Set("HX-Trigger", triggersJson)
		}

		// render
		return component.Render(rw.Request.Context(), rw.ResponseWriter)
	}

	// apply middlewares to the created handler
	// (they will execute before the handler at runtime)
	for _, m := range c.middlewares {
		handler = m(handler)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(&RW{
			Request:        r,
			ResponseWriter: w,
			target:         "",
		})
		if err != nil {
			panic(err)
		}
	}
}
