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

func (qcgh qrCodeGetHandler) ServeHTTP(reponseWriter http.ResponseWriter, request *http.Request) {
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
