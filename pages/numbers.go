package pages

import (
	"github.com/a-h/templ"
	gorouter "github.com/jmarren/go-router"
	"github.com/jmarren/go-router/views"
)

var NumbersPage *gorouter.ComponentRouter

func wrapNumbers(rw *gorouter.RW, component templ.Component) (templ.Component, error) {
	return views.NumbersNester(component), nil
}

func init() {
	NumbersPage = gorouter.CreateComponentRouter()
	NumbersPage.UseScripts("numbers.js").Trigger("numbers", "").PrefixWrap("/numbers", wrapNumbers)
	NumbersPage.GetComponent("/one", gorouter.SimpleComponent(views.One)).UseScripts("one.js").Trigger("hi", "")
	NumbersPage.GetComponent("/two", gorouter.SimpleComponent(views.Two)).UseScripts("two.js").Trigger("bye", "")
	NumbersPage.GetComponent("/say-hi", gorouter.SimpleComponent(views.Hi)).DontWrap()
}
