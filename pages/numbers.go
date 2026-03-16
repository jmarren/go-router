package pages

import (
	"github.com/a-h/templ"
	gorouter "github.com/jmarren/go-router"
	"github.com/jmarren/go-router/views"
)

func handleOne(rw *gorouter.RW) templ.Component {
	return views.One()
}

func handleTwo(rw *gorouter.RW) templ.Component {
	return views.Two()
}

var NumbersPage *gorouter.ComponentRouter

func wrapNumbers(rw *gorouter.RW, component templ.Component) (templ.Component, error) {
	return views.NumbersNester(component), nil
}

// func wrapNumbers(rw *gorouter.RW, component templ.Component) (templ.Component, error) {
// 	if strings.Contains(rw.CurrentUrl(), "numbers") {
// 		return component, nil
// 	}
//
// 	return views.NumbersNester(component), nil
// }

func init() {
	NumbersPage = gorouter.CreateComponentRouter()
	NumbersPage.UseScripts("numbers.js").Trigger("numbers", "")
	// NumbersPage.UseWrapFunc(wrapNumbers)
	NumbersPage.WrapWithoutSubpath("numbers", wrapNumbers)
	NumbersPage.GetComponent("/one", gorouter.UnsafeComponent(handleOne)).UseScripts("one.js").Trigger("hi", "")
	NumbersPage.GetComponent("/two", gorouter.UnsafeComponent(handleTwo)).UseScripts("two.js").Trigger("bye", "")

}
