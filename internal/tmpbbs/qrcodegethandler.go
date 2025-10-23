package tmpbbs

import (
	"net/http"

	"github.com/skip2/go-qrcode"
)

type qrCodeGetHandler struct{}

const qrCodeSize = 256

func newQRCodeGetHandler() *qrCodeGetHandler {
	return &qrCodeGetHandler{}
}

// ServeHTTP serves a QR code generated from a URL.
func (qcgh *qrCodeGetHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	url := request.URL.Query().Get("url")

	png, err := qrcode.Encode(url, qrcode.Medium, qrCodeSize)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Type", "image/png")

	if _, err = responseWriter.Write(png); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
}
