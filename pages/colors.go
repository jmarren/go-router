package pages

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	gorouter "github.com/jmarren/go-router"
	"github.com/jmarren/go-router/middleware"
	"github.com/jmarren/go-router/views"
)

func NestRed(rw *gorouter.RW, component templ.Component) templ.Component {
	if strings.Contains(rw.URL.Path, "red") {
		return views.RedNester(component)
	}
	return component
}

func catchRedError(rw *gorouter.RW, err error) (templ.Component, error) {
	fmt.Printf("caught red error: %s\n", err)
	return views.Red(), nil
}

var ColorsPage *gorouter.ComponentRouter

func init() {
	ColorsPage = gorouter.CreateComponentRouter()
	ColorsPage.UsePrefixWrap()
	ColorsPage.Retarget("#container")
	ColorsPage.Wrap(gorouter.SimpleWrapper(views.ColorsPage))
	ColorsPage.Use(middleware.LogUsernameMiddleware)
	ColorsPage.GetComponent("/red", gorouter.SimpleComponent(views.Red))
	ColorsPage.GetComponent("/yellow", gorouter.SimpleComponent(views.Yellow))
	ColorsPage.GetComponent("/green", gorouter.SimpleComponent(views.Red)).Catch(catchRedError)
}
