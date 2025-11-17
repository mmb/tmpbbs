package tmpbbs

import (
	"net/http"
)

type commonHeadersHandler struct {
	wrappedHandler http.Handler
}

func newCommonHeadersHandler(wrappedHandler http.Handler) *commonHeadersHandler {
	return &commonHeadersHandler{
		wrappedHandler: wrappedHandler,
	}
}

// ServeHTTP sets common headers then calls the wrapped handler.
func (chh *commonHeadersHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Cache-Control", "private, no-cache")
	responseWriter.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'; img-src *; "+
		"object-src 'none'; style-src *;")
	responseWriter.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	responseWriter.Header().Set("Referrer-Policy", "no-referrer")
	responseWriter.Header().Set("X-Content-Type-Options", "nosniff")

	chh.wrappedHandler.ServeHTTP(responseWriter, request)
}
