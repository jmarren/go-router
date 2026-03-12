package main

import (
	"fmt"
	"net/http"

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

func handleJohn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("john"))
}

func isHxRequest(w http.ResponseWriter, r *http.Request) bool {
	return r.Header.Get("HX-Request") != "true"
}

func handleHi(w http.ResponseWriter, r *http.Request) templ.Component {
	return views.Hi()
}

func main() {
	app := gorouter.CreateApp()

	// defines the outer component
	// app.UsePage(views.Base)

	// app.GetContent("/home",)

	// basePage := app.page(views.Base)
	// basePage.Use(logger)
	// basePage.Get("/home", handlers.Home)
	// basePage

	base := gorouter.CreateNestedComponentRouter(views.Page, isHxRequest)

	base.Get("/home", handleRoot)

	base.Component("/hi", handleHi)

	// base.Component("/hi",

	app.Page("/", "GET", handleRoot)
	app.Use(loggerOne)
	app.Use(loggerTwo)
	app.Get("/log", handleLog)
	usersRouter := gorouter.CreateRouter()
	usersRouter.Get("", handleUsers)
	johnRouter := gorouter.CreateRouter()
	johnRouter.Get("", handleJohn)
	usersRouter.SubRoute("/john", johnRouter)
	app.SubRoute("/users", usersRouter)
	app.Component("/johno", handleRoot).Get()

	app.Serve(":6060")

}
