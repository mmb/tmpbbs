package tmpbbs

import (
	"fmt"
	"net/http"
	"time"
)

// immutableGetHandler is an http.Handler that wraps requests for resources
// that are immutable for the life of the process or immutable for the same
// git commit. It handles conditional GETs and sets eTag and
// Last-Modified-Since headers.
type immutableGetHandler struct {
	wrappedHandler http.Handler
	startTime      time.Time
	eTag           string
	lastModified   string
}

func newImmutableGetHandler(wrappedHandler http.Handler) *immutableGetHandler {
	startTime := time.Now()

	handler := immutableGetHandler{
		wrappedHandler: wrappedHandler,
		startTime:      startTime.Truncate(time.Second).Round(0),
		lastModified:   startTime.Format(time.RFC1123),
	}
	if Commit != "" {
		handler.eTag = fmt.Sprintf("%q", Commit)
	}

	return &handler
}

func (igh *immutableGetHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if igh.eTag != "" {
		responseWriter.Header().Set("ETag", igh.eTag)
	}

	responseWriter.Header().Set("Last-Modified", igh.lastModified)

	if ifNoneMatch := request.Header.Get("If-None-Match"); igh.eTag != "" && ifNoneMatch == igh.eTag {
		responseWriter.WriteHeader(http.StatusNotModified)

		return
	}

	if ifModifiedSinceHeader := request.Header.Get("If-Modified-Since"); ifModifiedSinceHeader != "" {
		ifModifiedSince, err := time.Parse(time.RFC1123, ifModifiedSinceHeader)
		if err == nil && !igh.startTime.After(ifModifiedSince) {
			responseWriter.WriteHeader(http.StatusNotModified)

			return
		}
	}

	igh.wrappedHandler.ServeHTTP(responseWriter, request)
}
