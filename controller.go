package gorouter

type Controller struct {
	*ResponseWriter
	*Request
	// wrappers []Wrapper
}

// func (c *Controller) UseWrapper(w Wrapper) {
// 	c.wrappers = append([]Wrapper{w}, c.wrappers...)
// }

// type ControllerMiddleware func(c *Controller) *Controller
//
// func (c *Controller) Render() {
// 	c.ResponseWriter.Render(c.Context())
// }
//
// func WrapController(c *Controller) *Controller {
// 	c.component = views.Page(c.component, "hi")
// 	return c
// }
//
// // func (c *Controller) Serve() {
// //
// // }
