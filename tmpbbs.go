package main

import (
	"crypto/rand"
	"embed"
	"io/fs"
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
	pflag.StringSliceP("css-urls", "u", []string{"/static/main.css"}, "comma-separated list of CSS URLs ($TMPBBS_CSS_URLS)")
	pflag.BoolP("replies", "r", true, "Enable replies ($TMPBBS_REPLIES)")
	pflag.BoolP("emoji", "m", true, "Enable emoji shortcode expansion ($TMPBBS_EMOJI)")
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

//go:embed static
var staticFS embed.FS

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

	repliesPerPage := viper.GetInt("replies-per-page")
	postPostHandler := tmpbbs.NewPostPostHandler(repliesPerPage, postStore, tripCoder)
	repliesEnabled := viper.GetBool("replies")
	if repliesEnabled {
		http.Handle("POST /", postPostHandler)
		http.Handle("POST /{parentID}", postPostHandler)
	}
	postGetHandler := tmpbbs.NewPostGetHandler(title, repliesPerPage, viper.GetStringSlice("css-urls"), repliesEnabled, viper.GetBool("emoji"), postStore)
	http.Handle("GET /", postGetHandler)
	http.Handle("GET /{id}", postGetHandler)

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticDir)))
	http.Handle("GET /robots.txt", http.FileServerFS(staticDir))

	tlsCert := viper.GetString("tls-cert")
	tlsKey := viper.GetString("tls-key")
	listenAddress := viper.GetString("listen-address")
	if tlsCert != "" && tlsKey != "" {
		log.Fatal(http.ListenAndServeTLS(listenAddress, tlsCert, tlsKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(listenAddress, nil))
	}
}
