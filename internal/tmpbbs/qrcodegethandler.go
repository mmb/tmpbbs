package tmpbbs

import (
	"net/http"

	"github.com/skip2/go-qrcode"
)

const qrCodeSize = 256

// qrCodeGet serves a QR code generated from a URL.
func qrCodeGet(responseWriter http.ResponseWriter, request *http.Request) {
	png, err := qrcode.Encode(request.URL.Query().Get("url"), qrcode.Medium, qrCodeSize)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Type", "image/png")

	_, err = responseWriter.Write(png) // #nosec G705 -- image output
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

		return
	}
}
