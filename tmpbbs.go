package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mmb/tmpbbs/internal/tmpbbs"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.StringP("listen-address", "l", ":8080", "<host>:port to listen on ($TMPBBS_LISTEN_ADDRESS)")
	pflag.StringP("tls-cert", "c", "", "path to PEM server certificate ($TMPBBS_TLS_CERT)")
	pflag.StringP("tls-key", "k", "", "path to PEM server key ($TMPBBS_TLS_KEY)")
	pflag.StringP("title", "t", "tmpbbs", "site title ($TMPBBS_TITLE)")
	pflag.StringP("trip-code-salt", "a", "", "random salt to use for generating trip codes ($TMPBBS_TRIP_CODE_SALT)")
	pflag.StringP("load-posts", "p", "", `path to YAML or JSON file of posts to load, format [{"title":"","author":"","body":""}] ($TMPBBS_LOAD_POSTS)`)
	pflag.IntP("replies-per-page", "e", 10, "Number of replies to show per page ($TMPBBS_REPLIES_PER_PAGE)")
	pflag.StringSliceP("css-urls", "u", []string{"/css"}, "comma-separated list of CSS URLs ($TMPBBS_CSS_URLS)")
	pflag.BoolP("replies", "r", true, "Enable replies ($TMPBBS_REPLIES)")
	pflag.BoolP("help", "h", false, "usage help")

	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatal(err)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("tmpbbs")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func main() {
	if viper.GetBool("help") {
		pflag.CommandLine.SortFlags = false
		pflag.Usage()
		os.Exit(0)
	}

	title := viper.GetString("title")
	postStore, err := tmpbbs.NewPostStore(title)
	if err != nil {
		log.Fatal(err)
	}

	tripCoder, err := tmpbbs.NewTripCoder(viper.GetString("trip-code-salt"), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	loadPath := viper.GetString("load-posts")
	if loadPath != "" {
		err = postStore.LoadYAML(loadPath, tripCoder)
		if err != nil {
			log.Fatal(err)
		}
	}

	postPostHandler := tmpbbs.NewPostPostHandler(postStore, tripCoder)
	repliesEnabled := viper.GetBool("replies")
	if repliesEnabled {
		http.Handle("POST /", postPostHandler)
		http.Handle("POST /{parentID}", postPostHandler)
	}
	postGetHandler := tmpbbs.NewPostGetHandler(title, viper.GetInt("replies-per-page"), viper.GetStringSlice("css-urls"), repliesEnabled, postStore)
	http.Handle("GET /", postGetHandler)
	http.Handle("GET /{id}", postGetHandler)
	http.Handle("GET /css", new(tmpbbs.CSSHandler))
	http.Handle("GET /robots.txt", new(tmpbbs.RobotsHandler))

	tlsCert := viper.GetString("tls-cert")
	tlsKey := viper.GetString("tls-key")
	listenAddress := viper.GetString("listen-address")
	if tlsCert != "" && tlsKey != "" {
		log.Fatal(http.ListenAndServeTLS(listenAddress, tlsCert, tlsKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(listenAddress, nil))
	}
}
