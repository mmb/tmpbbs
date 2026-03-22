package tmpbbs

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Serve creates and configures an [http.Server] then starts listening.
func Serve(ctx context.Context, listenAddress string, tlsCertFile string, tlsKeyFile string,
	serveMux http.Handler,
) error {
	server := &http.Server{
		Addr:              listenAddress,
		Handler:           newLogHandler(serveMux),
		IdleTimeout:       120 * time.Second, //nolint:mnd // not worth creating a const
		ReadHeaderTimeout: 2 * time.Second,   //nolint:mnd // not worth creating a const
		ReadTimeout:       5 * time.Second,   //nolint:mnd // not worth creating a const
		WriteTimeout:      10 * time.Second,  //nolint:mnd // not worth creating a const
	}
	tlsEnabled := tlsCertFile != "" && tlsKeyFile != ""
	slog.InfoContext(ctx, "listening for HTTP", "address", listenAddress, "tlsEnabled", tlsEnabled)

	if tlsEnabled {
		return server.ListenAndServeTLS(tlsCertFile, tlsKeyFile)
	}

	return server.ListenAndServe()
}
