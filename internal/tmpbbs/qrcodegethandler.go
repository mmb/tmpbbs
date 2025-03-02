package tmpbbs

import (
	"net/http"

	"github.com/skip2/go-qrcode"
)

type qrCodeGetHandler struct{}

func NewQRCodeGetHandler() *qrCodeGetHandler {
	return &qrCodeGetHandler{}
}

func (qcgh qrCodeGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

	url := r.URL.Query().Get("url")

	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(png)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
