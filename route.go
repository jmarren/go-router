package gorouter

import "net/http"

type Route struct {
	path        string
	method      string
	middlewares []Middleware
	handler     func(w http.ResponseWriter, r *http.Request)
}

func (r *Route) HTTPHandler() http.HandlerFunc {
	// make a copy of the handler.
	// NOTE: we could use mutation instead, but cloning will prevent
	// issues if this method is called more than once
	handler := r.handler

	// apply middlewares to it
	for i := 0; i < len(r.middlewares); i++ {
		handler = r.middlewares[i](handler)
	}

	return handler

}
