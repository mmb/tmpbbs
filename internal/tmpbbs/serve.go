package tmpbbs

import (
	"net/http"
	"time"
)

func Serve(listenAddress string, tlsCertFile string, tlsKeyFile string) error {
	server := &http.Server{
		Addr:              listenAddress,
		IdleTimeout:       120 * time.Second, //nolint:mnd
		ReadHeaderTimeout: 2 * time.Second,   //nolint:mnd
		ReadTimeout:       5 * time.Second,   //nolint:mnd
		WriteTimeout:      10 * time.Second,  //nolint:mnd
	}

	if tlsCertFile != "" && tlsKeyFile != "" {
		return server.ListenAndServeTLS(tlsCertFile, tlsKeyFile)
	}

	return server.ListenAndServe()
}
