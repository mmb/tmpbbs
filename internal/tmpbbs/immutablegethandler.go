package tmpbbs

import (
	"net/http"
)

// immutableGetHandler is an [http.Handler] that sets the Cache-Control header
// to private, 1 year max age and immutable on success responses.
type immutableGetHandler struct {
	wrappedHandler http.Handler
}

type immutableResponseWriter struct {
	http.ResponseWriter

	wroteHeader bool
}

func newImmutableGetHandler(wrappedHandler http.Handler) *immutableGetHandler {
	return &immutableGetHandler{
		wrappedHandler: wrappedHandler,
	}
}

func (igh *immutableGetHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	irw := immutableResponseWriter{ResponseWriter: responseWriter}
	igh.wrappedHandler.ServeHTTP(&irw, request)
}

func (irw *immutableResponseWriter) Write(bytes []byte) (int, error) {
	if !irw.wroteHeader {
		irw.WriteHeader(http.StatusOK)
	}

	return irw.ResponseWriter.Write(bytes)
}

func (irw *immutableResponseWriter) WriteHeader(statusCode int) {
	if !irw.wroteHeader && statusCode >= http.StatusOK {
		irw.wroteHeader = true
		if statusCode == http.StatusOK || statusCode == http.StatusNotModified {
			irw.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
		}
	}

	irw.ResponseWriter.WriteHeader(statusCode)
}
