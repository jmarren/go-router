package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	gorouter "github.com/jmarren/go-router"
	"github.com/jmarren/go-router/middleware"
	"github.com/jmarren/go-router/pages"
	"github.com/jmarren/go-router/views"
)

func handleRoot(rw gorouter.RW) (templ.Component, error) {
	return views.Home(), nil
}

func SayHi(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("hi there"))
	return nil
}

func handleHi(rw gorouter.RW) (templ.Component, error) {
	return views.Hi(), nil
}

func Page(rw gorouter.RW, content templ.Component) (templ.Component, error) {

	fmt.Println("wrapping page!")

	usernameResponse := rw.Request.Context().Value("username")

	username, ok := usernameResponse.(string)

	if !ok {
		return content, fmt.Errorf("username not found")
	}

	return views.Page(content, username), nil
}

func pageCatcher(rw gorouter.RW, component templ.Component, err error) (templ.Component, error) {
	fmt.Printf("caught page err = %s\n", err)
	if err.Error() == "username not found" {
		return views.Page(component, "user not found"), nil
	}

	return component, err
}

func main() {
	app := gorouter.CreateApp()
	app.UseStaticDir("./static")

	app.UseBaseWrapper(views.Base)
	// app.UseScripts("/static/index.js")

	app.Use(middleware.Logger)
	// simple wrap the base component
	// app.SimpleHxWrap(views.Base)
	// hx-wrap the Page function and catch errors with the PageCatcher
	app.HxWrap(Page).Catch(pageCatcher)
	app.GetComponent("/", handleRoot)
	app.GetComponent("/hi", handleHi)

	// api := gorouter.CreateRouter()
	// api.Get("/hi", SayHi)
	// colorsPage.SubRoute("/api", api)

	app.SubComponent("/colors", pages.ColorsPage)
	app.SubComponent("/numbers", pages.NumbersPage)

	app.Serve(":6060")

}
