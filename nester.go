package gorouter

import (
	"net/http"

	"github.com/a-h/templ"
)

// any function that takes in a templ component and returns another
type SimpleNester func(c templ.Component) templ.Component

// any function that takes in a req, res, and component and returns a templ component
// this is used for wrapping subcomponents
type Nester func(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component

// // any function that takes in a req, res, and a slice of components and returns a single templ component
// type MultiNester func(w http.ResponseWriter, r *http.Request, components ...templ.Component) templ.Component

// converts a SimpleNester into a Nester that will always wrap
func FromSimple(fn SimpleNester) Nester {
	return func(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component {
		return fn(component)
	}
}

// converts a SimpleNester into a Nester
// that will wrap the component only if
// the request has HX-Request == true
func HxReqNester(fn SimpleNester) Nester {
	return func(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component {
		if r.Header.Get("HX-Request") == "true" {
			return component
		}
		return fn(component)
	}
}
