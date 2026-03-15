package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	gorouter "github.com/jmarren/go-router"
	"github.com/jmarren/go-router/views"
)

func loggerOne(h gorouter.Handler) gorouter.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println(r.Method + " " + r.URL.Path)
		return h(w, r)
	}
}

func loggerTwo(h gorouter.Handler) gorouter.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Println("logger 2")
		return h(w, r)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	return views.Home(), nil
}

func SayHi(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("hi there"))
	return nil
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("users"))
}

func handleYellow(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("yellow"))
	return fmt.Errorf("yellow is dumb")
}

func yellowCatcher(w http.ResponseWriter, r *http.Request, err error) error {
	fmt.Printf("caught yellow error = %s\n", err)
	return nil
}

func handleJohn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("john"))
}

func isHxRequest(w http.ResponseWriter, r *http.Request) bool {
	return r.Header.Get("HX-Request") != "true"
}

// func NestBase(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component {
// 	if isHxRequest(w, r) {
// 		return views.Page(component)
// 	}
// 	return component
// }

func NestColors(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component {
	return views.ColorsPage(component)
}

func handleHi(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	return views.Hi(), nil
}

func handleRed(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	return views.Red(), fmt.Errorf("some error")
}

func NestRed(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component {
	if strings.Contains(r.URL.Path, "red") {
		return views.RedNester(component)
	}
	return component
}

func NestNumbers(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component {
	return views.NumbersNester(component)
}

func handleOne(w http.ResponseWriter, r *http.Request) templ.Component {
	return views.One()
}

func handleTwo(w http.ResponseWriter, r *http.Request) templ.Component {
	return views.Two()
}

func catchRedError(w http.ResponseWriter, r *http.Request, err error) (templ.Component, error) {
	fmt.Printf("caught red error: %s\n", err)
	return views.Red(), nil
}

// add username to the request
func userMiddleware(h gorouter.Handler) gorouter.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		r = r.WithContext(context.WithValue(r.Context(), "username", "john"))
		return h(w, r)
	}
}

func logUsernameMiddleware(h gorouter.Handler) gorouter.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Printf("username = %s\n", r.Context().Value("username"))
		return h(w, r)
	}
}

func colorsCatcher(w http.ResponseWriter, r *http.Request, err error) (templ.Component, error) {
	fmt.Printf("uncaught colors error = %s\n", err)
	return views.DefaultErr(), nil
}

// type WrapperFunc func(w http.ResponseWriter, r *http.Request, component templ.Component) (templ.Component, error)
func Page(w http.ResponseWriter, r *http.Request, content templ.Component) (templ.Component, error) {

	fmt.Println("wrapping page!")

	usernameResponse := r.Context().Value("username")

	username, ok := usernameResponse.(string)

	if !ok {
		return content, fmt.Errorf("username not found")
	}

	return views.Page(content, username), nil
}

func PageCatcher(w http.ResponseWriter, r *http.Request, component templ.Component, err error) (templ.Component, error) {
	fmt.Printf("caught page err = %s\n", err)
	if err.Error() == "username not found" {
		return views.Page(component, "user not found"), nil
	}

	return component, err
}

func main() {
	app := gorouter.CreateApp()
	app.Use(loggerOne)
	app.HxWrap(Page).Catch(PageCatcher)
	app.GetComponent("/", handleRoot)
	app.GetComponent("/hi", handleHi)

	colorsPage := gorouter.CreateComponentRouter()
	colorsPage.SimpleHxWrap(views.ColorsPage)
	colorsPage.Use(logUsernameMiddleware)

	colorsPage.Use(loggerTwo)
	colorsPage.UseCatcher(yellowCatcher)
	// colorsPage.Use
	colorsPage.Get("/yellow", handleYellow)
	colorsPage.GetComponent("/red", handleRed)
	colorsPage.GetComponent("/green", handleRed).Catch(catchRedError)

	numbersPage := gorouter.CreateComponentRouter()
	numbersPage.SimpleHxWrap(views.ColorsPage)
	numbersPage.GetComponent("/one", gorouter.UnsafeComponent(handleOne))
	numbersPage.GetComponent("/two", gorouter.UnsafeComponent(handleTwo))

	api := gorouter.CreateRouter()
	api.Get("/hi", SayHi)
	colorsPage.SubRoute("/api", api)

	app.SubComponent("/colors", colorsPage)
	app.SubComponent("/numbers", numbersPage)

	app.Serve(":6060")

}
