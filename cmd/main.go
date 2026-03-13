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

func loggerOne(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method + " " + r.URL.Path)
		h(w, r)
	}
}

func loggerTwo(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("logger 2")
		h(w, r)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) templ.Component {
	return views.Home()
}

func SayHi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi there"))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("users"))
}

func handleYellow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("yellow"))
}

func handleJohn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("john"))
}

func isHxRequest(w http.ResponseWriter, r *http.Request) bool {
	return r.Header.Get("HX-Request") != "true"
}

func NestBase(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component {
	if isHxRequest(w, r) {
		return views.Page(component)
	}
	return component
}

func NestColors(w http.ResponseWriter, r *http.Request, component templ.Component) templ.Component {
	return views.ColorsPage(component)
}

func handleHi(w http.ResponseWriter, r *http.Request) templ.Component {
	return views.Hi()
}

func handleRed(w http.ResponseWriter, r *http.Request) templ.Component {
	return views.Red()
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

// add username to the request
func userMiddleware(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "username", "john"))
		h(w, r)
	}
}

func logUsernameMiddleware(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("username = %s\n", r.Context().Value("username"))
		h(w, r)
	}
}

func main() {
	app := gorouter.CreateApp()

	app.UseNester(NestBase)
	app.GetComponent("/", handleRoot)
	app.GetComponent("/hi", handleHi)
	app.Use(loggerOne)
	app.Use(userMiddleware)

	colorsPage := gorouter.CreateComponentRouter()
	colorsPage.Use(logUsernameMiddleware)
	colorsPage.UseNester(NestColors)
	colorsPage.Use(loggerTwo)
	colorsPage.Get("/yellow", handleYellow)
	colorsPage.GetComponent("/red", handleRed)
	colorsPage.GetComponent("/green", handleRed)

	numbersPage := gorouter.CreateComponentRouter()
	numbersPage.UseNester(NestNumbers)
	numbersPage.GetComponent("/one", handleOne)
	numbersPage.GetComponent("/two", handleTwo)

	api := gorouter.CreateRouter()
	api.Get("/hi", SayHi)
	colorsPage.SubRoute("/api", api)

	app.SubComponent("/colors", colorsPage)
	app.SubComponent("/numbers", numbersPage)

	app.Serve(":6060")

}
