package tmpbbs

import (
	"log/slog"
	"net/http"
	"time"
)

// Serve creates and configures an http.Server then starts listening.
func Serve(listenAddress string, tlsCertFile string, tlsKeyFile string, serveMux *http.ServeMux) error {
	server := &http.Server{
		Addr:              listenAddress,
		Handler:           newLogHandler(serveMux),
		IdleTimeout:       120 * time.Second, //nolint:mnd // not worth creating a const
		ReadHeaderTimeout: 2 * time.Second,   //nolint:mnd // not worth creating a const
		ReadTimeout:       5 * time.Second,   //nolint:mnd // not worth creating a const
		WriteTimeout:      10 * time.Second,  //nolint:mnd // not worth creating a const
	}

	if tlsCertFile != "" && tlsKeyFile != "" {
		slog.Info("listening for HTTP", "address", listenAddress, "tlsEnabled", true)

		return server.ListenAndServeTLS(tlsCertFile, tlsKeyFile)
	}

	slog.Info("listening for HTTP", "address", listenAddress, "tlsEnabled", false)

	return server.ListenAndServe()
}
