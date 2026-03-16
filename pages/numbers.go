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

func init() {
	NumbersPage = gorouter.CreateComponentRouter()
	/*
		wrapper = func (rw *RW, component templ.Component) (templ.Component, error) {
			return component, nil
		}
	*/

	NumbersPage.UsePrefixWrap()
	NumbersPage.Wrapper().UseFunc(wrapNumbers)

	/*
		func(rw *RW, component templ.Component) (templ.Component, error) {
			var err error
			component, err = func (rw *RW, component templ.Component) (templ.Component, error) {
				return component, nil
			}


			if err != nil {
				return component, err
			}

			return wrapNumbers(rw, component)
		}

		wrapper = 	*/
	NumbersPage.UseScripts("numbers.js").Trigger("numbers", "")
	NumbersPage.GetComponent("/one", gorouter.SimpleComponent(views.One)).UseScripts("one.js").Trigger("hi", "")
	NumbersPage.GetComponent("/two", gorouter.SimpleComponent(views.Two)).UseScripts("two.js").Trigger("bye", "")
	NumbersPage.GetComponent("/say-hi", gorouter.SimpleComponent(views.Hi)).DontWrap()

	// on mount

	/*
		func(rw *RW, component templ.Component) (templ.Component, error) {
			var err error
			component, err = func (rw *RW, component templ.Component) (templ.Component, error) {
				return component, nil
			}


			if err != nil {
				return component, err
			}

			return wrapNumbers(rw, component)
		}

		wrapper = 	*/

}
