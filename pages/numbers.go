package pages

import (
	"fmt"

	"github.com/a-h/templ"
	gorouter "github.com/jmarren/go-router"
	"github.com/jmarren/go-router/views"
)

var NumbersPage *gorouter.ComponentRouter

func wrapNumbers(rw *gorouter.RW, component templ.Component) (templ.Component, error) {
	fmt.Printf("wrapping numbers\n")
	return views.NumbersNester(component), nil
}

func handleOneEven(rw *gorouter.RW) (templ.Component, error) {
	fmt.Println("handling one even")
	return views.No(), nil
}

func handleTwoEven(rw *gorouter.RW) (templ.Component, error) {
	return views.Yes(), nil
}

func wrapIsEven(rw *gorouter.RW, component templ.Component) (templ.Component, error) {
	fmt.Printf("wrapping IsEven\n")
	return views.IsEven(component), nil
}

func init() {
	NumbersPage = gorouter.CreateComponentRouter()
	NumbersPage.UsePrefixWrap()
	NumbersPage.Wrap(wrapNumbers)
	NumbersPage.UseScripts("numbers.js").Trigger("numbers", "")
	NumbersPage.GetComponent("/one", gorouter.SimpleComponent(views.One)).UseScripts("one.js").Trigger("hi", "")
	NumbersPage.GetComponent("/two", gorouter.SimpleComponent(views.Two)).UseScripts("two.js").Trigger("bye", "")
	NumbersPage.GetComponent("/say-hi", gorouter.SimpleComponent(views.Hi)).DontWrap()

	SubNums := gorouter.CreateComponentRouter()

	SubNums.Wrap(wrapIsEven)
	SubNums.UsePrefixWrap()

	SubNums.GetComponent("/one", handleOneEven)
	SubNums.GetComponent("/two", handleTwoEven)

	NumbersPage.AddSubComponent("/is-even", SubNums)

}
