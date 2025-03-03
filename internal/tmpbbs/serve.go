package tmpbbs

import (
	"net/http"
	"time"
)

func Serve(listenAddress string, tlsCertFile string, tlsKeyFile string) error {
	server := &http.Server{
		Addr:              listenAddress,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	if tlsCertFile != "" && tlsKeyFile != "" {
		return server.ListenAndServeTLS(tlsCertFile, tlsKeyFile)
	}

	return server.ListenAndServe()
}
