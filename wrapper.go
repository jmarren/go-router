package gorouter

import (
	"github.com/a-h/templ"
)

type WrapperErrFunc func(rw *RW, component templ.Component, err error) (templ.Component, error)

type Wrapper interface {
	wrap(rw *RW, component templ.Component) (templ.Component, error)
	err(rw *RW, component templ.Component, err error) (templ.Component, error)
	Catch(errFunc WrapperErrFunc) Wrapper
	Use(m WrapFuncMiddleware) Wrapper
}

type WrapperFunc func(rw *RW, component templ.Component) (templ.Component, error)

type wrapper struct {
	wrapperFunc WrapperFunc
	errFunc     WrapperErrFunc
}

func (wr *wrapper) Use(m WrapFuncMiddleware) Wrapper {
	wr.wrapperFunc = m(wr.wrapperFunc)
	return wr
}

func (wr *wrapper) wrap(rw *RW, component templ.Component) (templ.Component, error) {
	return wr.wrapperFunc(rw, component)
}

func (wr *wrapper) err(rw *RW, component templ.Component, err error) (templ.Component, error) {
	return wr.errFunc(rw, component, err)
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

func unsafeErr(rw *RW, component templ.Component, err error) (templ.Component, error) {
	return component, err
}

func createWrapper(wrapperFunc WrapperFunc, errFunc WrapperErrFunc) Wrapper {

	if errFunc == nil {
		errFunc = unsafeErr
	}

	return &wrapper{
		wrapperFunc,
		errFunc,
	}
}

// any function that takes in a req, res, and component and returns a templ component
// this is used for wrapping subcomponents
// type Wrapper func(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error)

type WrapMiddleware func(w Wrapper) Wrapper

type WrapFuncMiddleware func(w WrapperFunc) WrapperFunc

// converts a SimpleWrapper into a Wrapper that returns a nil error
func FromSimple(s SimpleWrapper) Wrapper {
	// create default functions
	wrapperFunc := func(rw *RW, component templ.Component) (templ.Component, error) {
		return s(component), nil
	}

	errFunc := func(rw *RW, component templ.Component, err error) (templ.Component, error) {
		return component, err
	}
	return createWrapper(wrapperFunc, errFunc)
}

// // converts a SimpleWrapper into a WrapMiddleware
// func MiddlewareFromSimple(s SimpleWrapper) WrapMiddleware {
// 	return func(w Wrapper) Wrapper {
// 		return FromSimple(s)
// 	}
// }

// applies the wrapper only if not an hx-request
func hxWrapMiddleware(wr WrapperFunc) WrapperFunc {
	return func(rw *RW, component templ.Component) (templ.Component, error) {

		if rw.Request.Header.Get("HX-Request") == "true" {
			return component, nil
		}
		return wr(rw, component)
	}
}

// any function that takes in a templ component and returns another
type SimpleWrapper func(c templ.Component) templ.Component

// type ComponentWrapper interface {
// 	wrap(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error)
// 	err(w http.ResponseWriter, r *http.Request, err error) (templ.Component, error)
// 	UseWrapper(w Wrapper) ComponentWrapper
// }
//
// // concrete struct representing the ComponentWrapper
// type componentWrapper struct {
// 	wrappers    []Wrapper
// 	errCatchers []ComponentErrCatcher
// }

// // adds a catcher to the componentWrappers errCatchers
// func (c *componentWrapper) Catch(catcher ComponentErrCatcher) {
// 	c.errCatchers = append([]ComponentErrCatcher{catcher}, c.errCatchers...)
// }

// // adds a wrapper to the componentWrappers wrappers
// func (c *componentWrapper) UseWrapper(w Wrapper) ComponentWrapper {
// 	c.wrappers = append([]Wrapper{w}, c.wrappers...)
// 	return c
// }

// // iterates through wrappers and applies them
// func (c *componentWrapper) wrap(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error) {
// 	var err error
// 	for _, wrapper := range c.wrappers {
// 		component, err = wrapper(w, r, component)
// 	}
//
// 	return component, err
// }
//
// // iterates through catchers and applies them
// func (c *componentWrapper) err(w http.ResponseWriter, r *http.Request, err error) (templ.Component, error) {
//
// 	var component templ.Component
// 	for _, catcher := range c.errCatchers {
// 		component, err = catcher(w, r, err)
// 		if err == nil {
// 			break
// 		}
// 	}
// 	return component, err
// }

// any function that takes in a req, res, and a slice of components and returns a single templ component
// type MultiNester func(w http.ResponseWriter, r *http.Request, components ...templ.Component) templ.Component

// converts a SimpleNester into a Nester that will always wrap
// func FromSimple(fn SimpleWrapper) Wrapper {
// 	return func(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error) {
// 		return fn(component), nil
// 	}
// }

// func UnsafeHxReqWrapper(fn SimpleWrapper) ComponentWrapper {
// 	return &componentWrapper{
// 		wrappers:    []Wrapper{SimpleHxReqWrapper(fn)},
// 		errCatchers: []ComponentErrCatcher{},
// 	}
// }

// converts a SimpleNester into a Nester
// that will wrap the component only if
// the request has HX-Request == true
// func SimpleHxReqWrapper(fn SimpleWrapper) Wrapper {
// 	return func(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error) {
// 		if r.Header.Get("HX-Request") == "true" {
// 			return component, nil
// 		}
// 		return fn(component), nil
// 	}
// }
