package tmpbbs

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

const httpShutdownTimeout = 10 * time.Second

// Serve creates and configures an [http.Server] then starts listening.
func Serve(ctx context.Context, listenAddress string, tlsCertFile string, tlsKeyFile string,
	serveMux http.Handler,
) error {
	server := &http.Server{
		Addr:              listenAddress,
		Handler:           newLogHandler(serveMux),
		IdleTimeout:       120 * time.Second, //nolint:mnd // not worth creating a const
		ReadHeaderTimeout: 10 * time.Second,  //nolint:mnd // not worth creating a const
		ReadTimeout:       15 * time.Second,  //nolint:mnd // not worth creating a const
		WriteTimeout:      5 * time.Minute,   //nolint:mnd // not worth creating a const
	}
	tlsEnabled := tlsCertFile != "" && tlsKeyFile != ""
	slog.InfoContext(ctx, "listening for HTTP", "address", listenAddress, "tlsEnabled", tlsEnabled)

	errChan := make(chan error, 1)

	go func() {
		if tlsEnabled {
			errChan <- server.ListenAndServeTLS(tlsCertFile, tlsKeyFile)

			return
		}

		errChan <- server.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		slog.InfoContext(ctx, "stopping HTTP server")
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), httpShutdownTimeout)
	defer cancel()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		_ = server.Close()

		return err
	}

	return nil
}
