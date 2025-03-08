let dialog = document.getElementById("qr")

let img = document.createElement("img")
img.src = "/qr?url=" + encodeURIComponent(window.location.href)
img.loading = "lazy"
dialog.prepend(img)

let button = document.getElementById("show-qr")
button.addEventListener("click", () => { dialog.showModal() })
button.style.display = ""
