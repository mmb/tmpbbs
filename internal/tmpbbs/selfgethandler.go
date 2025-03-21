package tmpbbs

import (
	"net/http"
	"os"
)

type SelfGetHandler struct{}

func NewSelfGetHandler() *SelfGetHandler {
	return &SelfGetHandler{}
}

func (sgh SelfGetHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	binaryPath, err := os.Executable()
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Disposition", `attachment; filename="tmpbbs"`)
	http.ServeFile(responseWriter, request, binaryPath)
}
