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
	statusCode            int
}

func newLogHandler(wrappedHandler http.Handler) *logHandler {
	return &logHandler{
		wrappedHandler: wrappedHandler,
	}
}

func (lh *logHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	lrw := newLoggingResponseWriter(responseWriter)
	start := time.Now()

	lh.wrappedHandler.ServeHTTP(lrw, request)
	slog.Info("HTTP request",
		"remoteAddr", request.RemoteAddr,
		"method", request.Method,
		"path", request.URL.Path,
		"query", request.URL.Query(),
		"proto", request.Proto,
		"referer", request.Referer(),
		"userAgent", request.UserAgent(),
		"statusCode", cmp.Or(lrw.statusCode, http.StatusOK),
		"elapsed", time.Since(start),
	)
}

func newLoggingResponseWriter(wrappedResponseWriter http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		wrappedResponseWriter: wrappedResponseWriter,
	}
}

func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.wrappedResponseWriter.Header()
}

func (lrw *loggingResponseWriter) Write(bytes []byte) (int, error) {
	return lrw.wrappedResponseWriter.Write(bytes)
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.wrappedResponseWriter.WriteHeader(statusCode)
}
