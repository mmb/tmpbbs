package tmpbbs

import (
	"net/http"
)

// immutableGetHandler is an [http.Handler] that sets the Cache-Control header
// to private, 1 year max age and immutable.
type immutableGetHandler struct {
	wrappedHandler http.Handler
}

func newImmutableGetHandler(wrappedHandler http.Handler) *immutableGetHandler {
	return &immutableGetHandler{
		wrappedHandler: wrappedHandler,
	}
}

// ServeHTTP sets the Cache-Control header to private, 1 year max age and
// immutable.
func (igh *immutableGetHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	igh.wrappedHandler.ServeHTTP(responseWriter, request)
}
