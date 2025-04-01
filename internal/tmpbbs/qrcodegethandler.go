package tmpbbs

import (
	"net/http"

	"github.com/skip2/go-qrcode"
)

// A QRCodeGetHandler is an http.Handler that serves QR code images.
type QRCodeGetHandler struct{}

const qrCodeSize = 256

// NewQRCodeGetHandler returns a new QRCodeGetHandler.
func NewQRCodeGetHandler() *QRCodeGetHandler {
	return &QRCodeGetHandler{}
}

func (qcgh QRCodeGetHandler) ServeHTTP(reponseWriter http.ResponseWriter, request *http.Request) {
	url := request.URL.Query().Get("url")

	png, err := qrcode.Encode(url, qrcode.Medium, qrCodeSize)
	if err != nil {
		http.Error(reponseWriter, err.Error(), http.StatusInternalServerError)

		return
	}

	reponseWriter.Header().Set("Content-Type", "image/png")

	if _, err = reponseWriter.Write(png); err != nil {
		http.Error(reponseWriter, err.Error(), http.StatusInternalServerError)
	}
}
