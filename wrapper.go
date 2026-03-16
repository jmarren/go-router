package gorouter

import (
	"fmt"

	"github.com/a-h/templ"
)

type WrapperErrFunc func(rw *RW, component templ.Component, err error) (templ.Component, error)

type WrapperFunc func(rw *RW, component templ.Component) (templ.Component, error)

type WrapMiddleware func(w WrapperFunc) WrapperFunc

type Wrapper interface {
	wrapperFunc() func(rw *RW, component templ.Component) (templ.Component, error)
	Catch(errFunc WrapperErrFunc) Wrapper
	Use(m WrapMiddleware) Wrapper
	UseFunc(w WrapperFunc) Wrapper
}

type wrapper struct {
	wrapFunc WrapperFunc
	errFunc  WrapperErrFunc
}

func (wr *wrapper) wrapperFunc() func(rw *RW, component templ.Component) (templ.Component, error) {

	return func(rw *RW, component templ.Component) (templ.Component, error) {
		var err error
		component, err = wr.wrapFunc(rw, component)
		if err != nil {
			return wr.errFunc(rw, component, err)
		}

		return component, err
	}
}

func (wr *wrapper) err(rw *RW, component templ.Component, err error) (templ.Component, error) {
	return wr.errFunc(rw, component, err)
}

func (wr *wrapper) Use(fn WrapMiddleware) Wrapper {
	wr.wrapFunc = fn(wr.wrapFunc)
	return wr
}

func (wr *wrapper) Catch(errFunc WrapperErrFunc) Wrapper {

	// store the current errFunc of the wrapper
	curr := wr.errFunc

	// update the errFunc to try using the new errFunc first
	wr.errFunc = func(rw *RW, component templ.Component, err error) (templ.Component, error) {
		// use the new errFunc
		component, err = errFunc(rw, component, err)
		if err != nil {
			return curr(rw, component, err)
		}
		return component, err
	}

	return wr
}

func (wr *wrapper) UseFunc(w WrapperFunc) Wrapper {
	wr.Use(MiddlewareFromFunc(w))
	return wr
}

func PrefixWrap(prefix string) WrapMiddleware {
	return func(w WrapperFunc) WrapperFunc {
		return func(rw *RW, component templ.Component) (templ.Component, error) {
			fmt.Printf("checking if path %s has prefix %s\n", rw.URL.Path, prefix)
			if rw.PathHasPrefix(prefix) {
				fmt.Printf("path %s has prefix %s\n", rw.URL.Path, prefix)
				return component, nil
			}
			return w(rw, component)
		}
	}
}

func unsafeErr(rw *RW, component templ.Component, err error) (templ.Component, error) {
	return component, err
}

func defaultWrapFunc(rw *RW, component templ.Component) (templ.Component, error) {
	return component, nil
}

func defaultWrapper() Wrapper {
	return &wrapper{
		wrapFunc: defaultWrapFunc,
		errFunc:  unsafeErr,
	}
}

func MiddlewareFromFunc(wrapper WrapperFunc) WrapMiddleware {
	return func(w WrapperFunc) WrapperFunc {
		return func(rw *RW, component templ.Component) (templ.Component, error) {
			var err error
			component, err = w(rw, component)

			if err != nil {
				return component, err
			}

			return wrapper(rw, component)
		}
	}
}

func SimpleWrapper(s func(c templ.Component) templ.Component) WrapperFunc {
	return func(rw *RW, component templ.Component) (templ.Component, error) {
		return s(component), nil
	}
}

// converts a SimpleWrapper into a Wrapper that returns a nil error
func FromSimple(s func(c templ.Component) templ.Component) Wrapper {
	// create default functions
	wrapperFunc := func(rw *RW, component templ.Component) (templ.Component, error) {
		return s(component), nil
	}

	errFunc := func(rw *RW, component templ.Component, err error) (templ.Component, error) {
		return component, err
	}
	return &wrapper{
		wrapFunc: wrapperFunc,
		errFunc:  errFunc,
	}
}

// applies the wrapper only if not an hx-request
func hxWrapMiddleware(wr WrapperFunc) WrapperFunc {
	return func(rw *RW, component templ.Component) (templ.Component, error) {
		if rw.Request.Header.Get("HX-Request") == "true" {
			return component, nil
		}
		return wr(rw, component)
	}
}
