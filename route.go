package gorouter

import (
	"fmt"
	"net/http"
)

type RW struct {
	http.ResponseWriter
	*http.Request
}

type ErrCatcher func(rw RW, err error) error

type Handler func(rw RW) error

type Route struct {
	path        string
	method      string
	middlewares []Middleware
	handler     Handler
	catchers    []ErrCatcher
}

func (r *Route) Catch(c ErrCatcher) *Route {
	r.catchers = append([]ErrCatcher{c}, r.catchers...)
	return r
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

	return func(w http.ResponseWriter, req *http.Request) {

		rw := RW{
			ResponseWriter: w,
			Request:        req,
		}

		err := handler(rw)

		fmt.Printf("num catchers = %d\n", len(r.catchers))

		if err != nil {
			for _, catcher := range r.catchers {
				err = catcher(rw, err)
				if err == nil {
					break
				}
			}
			// panic if err is still not resolved
			if err != nil {
				panic(err)
			}
		}
	}

}
