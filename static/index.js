console.log('hi!')

var executed = []

htmx.on("htmx:configRequest", (event) => {
	event.detail.headers["HX-Executed"] = JSON.stringify(executed)
})

htmx.on("hi", (event) => {
	console.log("hi triggered")
})


htmx.on("bye", (event) => {
	console.log("bye triggered")
})


htmx.on("numbers", (event) => {
	console.log("numbers triggered")
})
