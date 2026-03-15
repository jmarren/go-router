package pages

import (
	"github.com/a-h/templ"
	gorouter "github.com/jmarren/go-router"
	"github.com/jmarren/go-router/views"
)

func handleOne(rw gorouter.RW) templ.Component {
	return views.One()
}

func handleTwo(rw gorouter.RW) templ.Component {
	return views.Two()
}

var NumbersPage *gorouter.ComponentRouter

func init() {
	NumbersPage = gorouter.CreateComponentRouter()
	NumbersPage.SimpleHxWrap(views.NumbersNester)
	NumbersPage.GetComponent("/one", gorouter.UnsafeComponent(handleOne))
	NumbersPage.GetComponent("/two", gorouter.UnsafeComponent(handleTwo))

}
