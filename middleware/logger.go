package middleware

import (
	"fmt"

	gorouter "github.com/jmarren/go-router"
)

func Logger(h gorouter.Handler) gorouter.Handler {
	return func(rw gorouter.RW) error {
		fmt.Println(rw.Method + " " + rw.URL.Path)
		return h(rw)
	}
}
