package main

import (
	"fmt"
	"net/http"

	gorouter "github.com/jmarren/go-router"
)

func main() {
	app := gorouter.CreateApp()

	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi there"))
	})

	app.Use(func(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Method + " " + r.URL.Path)
			h(w, r)
		}
	})

	app.Get("/log", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi there"))
	})

	usersRouter := gorouter.CreateRouter()

	usersRouter.Get("", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("users"))
	})

	app.SubRoute("/users", usersRouter)

	app.Serve(":6060")

}
