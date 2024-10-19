package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/mmb/tmpbbs/internal/tmpbbs"
)

func main() {
	var (
		cssURL        = flag.String("cssURL", "/css", "CSS URL")
		listenAddress = flag.String("listenAddress", ":8080", "<host>:port to listen on")
		title         = flag.String("title", "tmpbbs", "Site title")
		tlsCert       = flag.String("tlsCert", "", "path to PEM server certificate")
		tlsKey        = flag.String("tlsKey", "", "path to PEM server key")
		tripCodeSalt  = flag.String("tripCodeSalt", "", "random salt to use for generating trip codes")
	)
	flag.Parse()

	postStore, err := tmpbbs.NewPostStore(*title)
	if err != nil {
		log.Fatal(err)
	}

	tripCoder := tmpbbs.NewTripCoder(*tripCodeSalt)

	postPostHandler := tmpbbs.CreatePostPostHandler(postStore, tripCoder)
	http.HandleFunc("POST /", postPostHandler)
	http.HandleFunc("POST /{parentID}", postPostHandler)
	getPostHandler := tmpbbs.CreateGetPostHandler(postStore, cssURL, title)
	http.HandleFunc("GET /", getPostHandler)
	http.HandleFunc("GET /{id}", getPostHandler)
	http.HandleFunc("GET /css", tmpbbs.CSSHandler)
	http.HandleFunc("GET /robots.txt", tmpbbs.RobotsHandler)

	if *tlsCert != "" && *tlsKey != "" {
		log.Fatal(http.ListenAndServeTLS(*listenAddress, *tlsCert, *tlsKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(*listenAddress, nil))
	}
}
