package pages

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	gorouter "github.com/jmarren/go-router"
	"github.com/jmarren/go-router/middleware"
	"github.com/jmarren/go-router/views"
)

func handleRed(rw gorouter.RW) (templ.Component, error) {
	return views.Red(), fmt.Errorf("some error")
}

func NestRed(rw gorouter.RW, component templ.Component) templ.Component {
	if strings.Contains(rw.URL.Path, "red") {
		return views.RedNester(component)
	}
	return component
}

func catchRedError(rw gorouter.RW, err error) (templ.Component, error) {
	fmt.Printf("caught red error: %s\n", err)
	return views.Red(), nil
}

func handleYellow(rw gorouter.RW) error {
	rw.ResponseWriter.Write([]byte("yellow"))
	return fmt.Errorf("yellow is dumb")
}

func yellowCatcher(rw gorouter.RW, err error) error {
	fmt.Printf("caught yellow error = %s\n", err)
	return nil
}

var ColorsPage *gorouter.ComponentRouter

func init() {

	ColorsPage = gorouter.CreateComponentRouter()
	ColorsPage.SimpleHxWrap(views.ColorsPage)
	ColorsPage.Use(middleware.LogUsernameMiddleware)

	// colorsPage.Use(middleware.LoggerTwo)
	ColorsPage.UseCatcher(yellowCatcher)
	// colorsPage.Use
	ColorsPage.Get("/yellow", handleYellow)
	ColorsPage.GetComponent("/red", handleRed)
	ColorsPage.GetComponent("/green", handleRed).Catch(catchRedError)
}
