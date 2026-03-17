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

func handleRoot(rw *gorouter.RW) (templ.Component, error) {
	return views.Home(), nil
}

func SayHi(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("hi there"))
	return nil
}

func handleHi(rw *gorouter.RW) (templ.Component, error) {
	return views.Hi(), nil
}

func Page(rw *gorouter.RW, content templ.Component) (templ.Component, error) {

	fmt.Println("wrapping page!")

	usernameResponse := rw.Request.Context().Value("username")

	username, ok := usernameResponse.(string)

	if !ok {
		return content, fmt.Errorf("username not found")
	}

	return views.Page(content, username), nil
}

func pageCatcher(rw *gorouter.RW, component templ.Component, err error) (templ.Component, error) {
	fmt.Printf("caught page err = %s\n", err)
	if err.Error() == "username not found" {
		return views.Page(component, "user not found"), nil
	}

	return component, err
}

func PageWrap(w gorouter.WrapperFunc) gorouter.WrapperFunc {
	return func(rw *gorouter.RW, component templ.Component) (templ.Component, error) {
		var err error
		component, err = w(rw, component)
		if err != nil {
			return component, err
		}
		return Page(rw, component)
	}
}

func main() {
	app := gorouter.CreateApp()
	app.UseStaticDir("./static")
	app.Use(middleware.UserMiddleware)

	app.UseBaseWrapper(views.Base)
	app.Use(middleware.Logger)

	app.Wrap(Page).Catch(pageCatcher)
	app.HxWrap()
	app.GetComponent("/", handleRoot)
	app.GetComponent("/hi", handleHi)

	app.AddSubComponent("/colors", pages.ColorsPage)
	app.AddSubComponent("/numbers", pages.NumbersPage)

	app.Serve(":6060")

}

// api := gorouter.CreateRouter()
// api.Get("/hi", SayHi)
// colorsPage.SubRoute("/api", api)
