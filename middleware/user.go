package middleware

import (
	"context"
	"fmt"

	gorouter "github.com/jmarren/go-router"
)

// add username to the request
func UserMiddleware(h gorouter.Handler) gorouter.Handler {
	return func(rw *gorouter.RW) error {
		rw.Request = rw.Request.WithContext(context.WithValue(rw.Context(), "username", "john"))
		return h(rw)
	}
}

func LogUsernameMiddleware(h gorouter.Handler) gorouter.Handler {
	return func(rw *gorouter.RW) error {
		fmt.Printf("username = %s\n", rw.Context().Value("username"))
		return h(rw)
	}
}
