package gorouter

import "net/http"

type Request struct {
	*http.Request
}

func (r *Request) IsHxRequest() bool {
	return r.Header.Get("HX-Request") == "true"
}
