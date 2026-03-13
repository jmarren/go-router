package gorouter

import "net/http"

type Route struct {
	path        string
	method      string
	middlewares []Middleware
	handler     func(w http.ResponseWriter, r *http.Request)
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(route.middlewares); i++ {
		route.handler = route.middlewares[i](route.handler)
	}
	route.handler(w, r)
}
