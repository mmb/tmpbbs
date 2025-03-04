package tmpbbs

import (
	"net/http"

	"github.com/skip2/go-qrcode"
)

const qrCodeSize = 256

type QRCodeGetHandler struct{}

func NewQRCodeGetHandler() *QRCodeGetHandler {
	return &QRCodeGetHandler{}
}

func (qcgh QRCodeGetHandler) ServeHTTP(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Cache-Control", "no-store")

	url := request.URL.Query().Get("url")

	png, err := qrcode.Encode(url, qrcode.Medium, qrCodeSize)
	if err != nil {
		http.Error(reponseWriter, err.Error(), http.StatusInternalServerError)
	}

	reponseWriter.Header().Set("Content-Type", "image/png")

	if _, err = reponseWriter.Write(png); err != nil {
		http.Error(reponseWriter, err.Error(), http.StatusInternalServerError)
	}
}
