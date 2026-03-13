package main

import (
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

func handleLog(w http.ResponseWriter, r *http.Request) {
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

func main() {
	app := gorouter.CreateApp()

	app.UseNester(NestBase)
	app.GetComponent("/", handleRoot)
	app.GetComponent("/hi", handleHi)

	subComponent := gorouter.CreateComponentRouter()

	// subComponent := app.CreateSubRouter("/colors")

	subComponent.UseNester(NestRed)

	subComponent.Get("/yellow", handleYellow)
	subComponent.GetComponent("/red", handleRed)
	subComponent.GetComponent("/green", handleRed)

	app.SubComponent("/colors", subComponent)

	// app.SubRoute("/red", subComponent)

	// app.SubRoute("/red", subComponent)

	// defines the outer component
	// app.UsePage(views.Base)

	// app.GetContent("/home",)

	// basePage := app.page(views.Base)
	// basePage.Use(logger)
	// basePage.Get("/home", handlers.Home)
	// basePage

	// base := gorouter.CreateNestedComponentRouter(views.Page, isHxRequest)
	//
	// base.Get("/home", handleRoot)
	//
	// base.Component("/hi", handleHi)
	//
	// base.Component("/hi",

	// app.Page("/", "GET", handleRoot)
	app.Use(loggerOne)
	app.Use(loggerTwo)
	app.Get("/log", handleLog)
	usersRouter := gorouter.CreateRouter()
	usersRouter.Get("", handleUsers)
	johnRouter := gorouter.CreateRouter()
	johnRouter.Get("", handleJohn)
	// usersRouter.SubRoute("/john", johnRouter)
	// app.SubRoute("/users", usersRouter)
	// app.Component("/johno", handleRoot).Get()

	app.Serve(":6060")

}
