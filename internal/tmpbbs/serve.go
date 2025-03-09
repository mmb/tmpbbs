package tmpbbs

import (
	"net/http"
	"time"
)

func Serve(listenAddress string, tlsCertFile string, tlsKeyFile string, serveMux *http.ServeMux) error {
	server := &http.Server{
		Addr:              listenAddress,
		Handler:           serveMux,
		IdleTimeout:       120 * time.Second, //nolint:mnd // not worth creating a const
		ReadHeaderTimeout: 2 * time.Second,   //nolint:mnd // not worth creating a const
		ReadTimeout:       5 * time.Second,   //nolint:mnd // not worth creating a const
		WriteTimeout:      10 * time.Second,  //nolint:mnd // not worth creating a const
	}

	if tlsCertFile != "" && tlsKeyFile != "" {
		return server.ListenAndServeTLS(tlsCertFile, tlsKeyFile)
	}

	return server.ListenAndServe()
}
