package tmpbbs

import (
	"net/http"

	"github.com/skip2/go-qrcode"
)

type QRCodeGetHandler struct{}

func NewQRCodeGetHandler() *QRCodeGetHandler {
	return &QRCodeGetHandler{}
}

func (qcgh QRCodeGetHandler) ServeHTTP(reponseWriter http.ResponseWriter, request *http.Request) {
	reponseWriter.Header().Set("Cache-Control", "no-store")

	url := request.URL.Query().Get("url")

	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		http.Error(reponseWriter, err.Error(), http.StatusInternalServerError)
	}
	reponseWriter.Header().Set("Content-Type", "image/png")
	_, err = reponseWriter.Write(png)
	if err != nil {
		http.Error(reponseWriter, err.Error(), http.StatusInternalServerError)
	}
}
