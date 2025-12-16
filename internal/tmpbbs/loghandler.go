package tmpbbs

import (
	"cmp"
	"log/slog"
	"net/http"
	"time"
)

type logHandler struct {
	wrappedHandler http.Handler
}

type loggingResponseWriter struct {
	wrappedResponseWriter http.ResponseWriter
	responseSize          int
	statusCode            int
}

func newLogHandler(wrappedHandler http.Handler) *logHandler {
	return &logHandler{
		wrappedHandler: wrappedHandler,
	}
}

// ServeHTTP calls the wrapped handler then logs information about the request
// and response.
func (lh *logHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	lrw := newLoggingResponseWriter(responseWriter)
	start := time.Now()

	lh.wrappedHandler.ServeHTTP(lrw, request)
	slog.InfoContext(request.Context(), "HTTP request",
		slog.Group("request",
			"remoteAddr", request.RemoteAddr,
			"method", request.Method,
			"path", request.URL.Path,
			"query", request.URL.Query(),
			"proto", request.Proto,
			"referer", request.Referer(),
			"userAgent", request.UserAgent(),
		),
		slog.Group("response",
			"statusCode", cmp.Or(lrw.statusCode, http.StatusOK),
			"size", lrw.responseSize,
			"elapsed", time.Since(start),
		),
	)
}

func newLoggingResponseWriter(wrappedResponseWriter http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		wrappedResponseWriter: wrappedResponseWriter,
	}
}

// Header returns what the wrapped ResponseWriter's Header method returns.
func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.wrappedResponseWriter.Header()
}

// Write writes bytes to the wrapped ResponseWriter and increments the total
// bytes written.
func (lrw *loggingResponseWriter) Write(bytes []byte) (int, error) {
	bytesWritten, err := lrw.wrappedResponseWriter.Write(bytes)
	lrw.responseSize += bytesWritten

	return bytesWritten, err
}

// WriteHeader sets the status code on the wrapped ResponseWriter and calls its
// WriteHeader method.
func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.wrappedResponseWriter.WriteHeader(statusCode)
}
