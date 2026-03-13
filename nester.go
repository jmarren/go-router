package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
)

// any function that takes in a templ component and returns another
type SimpleNester func(c templ.Component) templ.Component

type ComponentWrapper interface {
	wrap(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error)
	err(w http.ResponseWriter, r *http.Request, err error) (templ.Component, error)
}

// any function that takes in a req, res, and component and returns a templ component
// this is used for wrapping subcomponents
type Wrapper func(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error)

// type Wrapper func(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error)
type ErrWrapper func(w http.ResponseWriter, r *http.Request, err error) (templ.Component, error)

type componentWrapper struct {
	wrapper    Wrapper
	errWrapper ErrWrapper
}

func (c *componentWrapper) wrap(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error) {
	return c.wrapper(w, r, component)
}

func (c *componentWrapper) err(w http.ResponseWriter, r *http.Request, err error) (templ.Component, error) {
	return c.errWrapper(w, r, err)
}

func newComponentWrapper(wrapper Wrapper, errWrapper ErrWrapper) *componentWrapper {
	return &componentWrapper{
		wrapper,
		errWrapper,
	}
}

// any function that takes in a req, res, and a slice of components and returns a single templ component
// type MultiNester func(w http.ResponseWriter, r *http.Request, components ...templ.Component) templ.Component

// converts a SimpleNester into a Nester that will always wrap
func FromSimple(fn SimpleNester) Wrapper {
	return func(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error) {
		return fn(component), nil
	}
}

func UnsafeHxReqWrapper(fn SimpleNester) ComponentWrapper {
	return &componentWrapper{
		wrapper:    HxReqWrapper(fn),
		errWrapper: nil,
	}
}

// converts a SimpleNester into a Nester
// that will wrap the component only if
// the request has HX-Request == true
func HxReqWrapper(fn SimpleNester) Wrapper {
	return func(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error) {
		if r.Header.Get("HX-Request") == "true" {
			return component, nil
		}
		return fn(component), nil
	}
}
