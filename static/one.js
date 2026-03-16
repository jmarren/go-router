console.log("one!")

function makeNumberGreen() {
	document.querySelector(".number").style.backgroundColor = "green"
}

makeNumberGreen()


htmx.on("htmx:load", (e) => { 
	if (e.detail.elt.id == "one") {
		makeNumberGreen()
	}
})

executed.push("one.js")
