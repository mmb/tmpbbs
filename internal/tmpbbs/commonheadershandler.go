package tmpbbs

import (
	"net/http"
)

type commonHeadersHandler struct {
	wrappedHandler http.Handler
	cspHeader      string
}

func newCommonHeadersHandler(wrappedHandler http.Handler, externalCSS bool) *commonHeadersHandler {
	cspHeader := "default-src 'self'; " +
		"base-uri 'self'; " +
		"form-action 'self'; " +
		"frame-ancestors 'none'; " +
		"img-src *; " +
		"object-src 'none'"
	if externalCSS {
		cspHeader += "; style-src *"
	}

	return &commonHeadersHandler{
		cspHeader:      cspHeader,
		wrappedHandler: wrappedHandler,
	}
}

// ServeHTTP sets common headers then calls the wrapped handler.
func (chh *commonHeadersHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Cache-Control", "private, no-cache")
	responseWriter.Header().Set("Content-Security-Policy", chh.cspHeader)
	responseWriter.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	responseWriter.Header().Set("Referrer-Policy", "no-referrer")
	responseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	responseWriter.Header().Set("X-Frame-Options", "DENY")
	chh.wrappedHandler.ServeHTTP(responseWriter, request)
}
