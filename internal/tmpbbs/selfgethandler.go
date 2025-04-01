package tmpbbs

import (
	"net/http"
	"os"
)

// A SelfGetHandler is an http.Handler that serves the running server binary.
type SelfGetHandler struct{}

// NewSelfGetHandler returns a new SelfGetHandler.
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
