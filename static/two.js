
console.log("two!")

function makeNumberBlue() {
	document.querySelector(".number").style.backgroundColor = "blue"
}

makeNumberBlue()

htmx.on("htmx:load", (e) => { 
	if (e.detail.elt.id == "two") {
		makeNumberBlue()
	}
})

executed.push("two.js")
