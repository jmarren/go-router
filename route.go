package gorouter

import "net/http"

type Route struct {
	path        string
	method      string
	middlewares []Middleware
	handler     func(w http.ResponseWriter, r *http.Request)
}

func (r *Route) HTTPHandler() http.HandlerFunc {
	// make a copy of the handler
	handler := r.handler

	// apply middlewares to it
	for i := 0; i < len(r.middlewares); i++ {
		handler = r.middlewares[i](handler)
	}

	return handler

}
