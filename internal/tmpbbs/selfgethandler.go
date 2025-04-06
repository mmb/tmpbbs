package tmpbbs

import (
	"net/http"
	"os"
)

type selfGetHandler struct{}

func newSelfGetHandler() *selfGetHandler {
	return &selfGetHandler{}
}

func (sgh selfGetHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	binaryPath, err := os.Executable()
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Disposition", `attachment; filename="tmpbbs"`)
	http.ServeFile(responseWriter, request, binaryPath)
}
