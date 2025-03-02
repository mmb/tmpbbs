package main

import (
	"crypto/rand"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mmb/tmpbbs/internal/tmpbbs"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var version = "unspecified"

//go:embed static
var staticFS embed.FS

func init() {
	pflag.StringP("listen-address", "l", ":8080", "<host>:port to listen on ($TMPBBS_LISTEN_ADDRESS)")
	pflag.StringP("tls-cert", "c", "", "path to PEM server certificate ($TMPBBS_TLS_CERT)")
	pflag.StringP("tls-key", "k", "", "path to PEM server key ($TMPBBS_TLS_KEY)")
	pflag.StringP("title", "t", "tmpbbs", "site title ($TMPBBS_TITLE)")
	pflag.StringP("trip-code-salt", "a", "", "random salt to use for generating trip codes ($TMPBBS_TRIP_CODE_SALT)")
	pflag.StringSliceP("load-posts", "p", []string{},
		`comma-separated paths of YAML or JSON files of posts to load, format [{"title":"","author":"","body":""}]
($TMPBBS_LOAD_POSTS)`)
	pflag.StringSliceP("serve-fs-paths", "f", []string{},
		"comma-separated list of urlprefix=/local/dir to serve ($TMPBBS_SERVE_FS_PATHS)")
	pflag.IntP("replies-per-page", "e", 10, "number of replies to show per page ($TMPBBS_REPLIES_PER_PAGE)")
	pflag.StringSliceP("css-urls", "u", []string{"/static/main.css"},
		"comma-separated list of CSS URLs ($TMPBBS_CSS_URLS)")
	pflag.BoolP("replies", "r", true, "enable replies ($TMPBBS_REPLIES)")
	pflag.BoolP("emoji", "m", true, "enable emoji shortcode expansion ($TMPBBS_EMOJI)")
	pflag.BoolP("qr-codes", "q", true, "enable shareable URL QR codes ($TMPBBS_QR_CODES)")
	pflag.BoolP("version", "v", false, "print version")
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

	if viper.GetBool("version") {
		fmt.Println(version)
		os.Exit(0)
	}

	postStore := tmpbbs.NewPostStore(viper.GetString("title"))

	tripCoder, err := tmpbbs.NewTripCoder(viper.GetString("trip-code-salt"), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	for _, loadPath := range viper.GetStringSlice("load-posts") {
		err = postStore.LoadYAML(loadPath, tripCoder)
		if err != nil {
			log.Fatal(err)
		}
	}

	repliesPerPage := viper.GetInt("replies-per-page")
	postPostHandler := tmpbbs.NewPostPostHandler(repliesPerPage, postStore, tripCoder)
	repliesEnabled := viper.GetBool("replies")

	if repliesEnabled {
		http.Handle("POST /{$}", postPostHandler)
		http.Handle("POST /{parentID}", postPostHandler)
	}

	qrCodesEnabled := viper.GetBool("qr-codes")
	if qrCodesEnabled {
		http.Handle("GET /qr", tmpbbs.NewQRCodeGetHandler())
	}

	postGetHandler := tmpbbs.NewPostGetHandler(repliesPerPage, viper.GetStringSlice("css-urls"), repliesEnabled,
		viper.GetBool("emoji"), qrCodesEnabled, postStore)
	http.Handle("GET /{$}", postGetHandler)
	http.Handle("GET /{id}", postGetHandler)

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("GET /static/", http.StripPrefix("/static", http.FileServerFS(staticDir)))
	http.Handle("GET /robots.txt", http.FileServerFS(staticDir))

	for _, dirMapping := range viper.GetStringSlice("serve-fs-paths") {
		parts := strings.SplitN(dirMapping, "=", 2)
		urlPrefix, dir := "/"+parts[0], parts[1]
		http.Handle(fmt.Sprintf("GET %s/", urlPrefix), http.StripPrefix(urlPrefix, http.FileServer(http.Dir(dir))))
	}

	tlsCert := viper.GetString("tls-cert")
	tlsKey := viper.GetString("tls-key")
	listenAddress := viper.GetString("listen-address")

	if tlsCert != "" && tlsKey != "" {
		log.Fatal(http.ListenAndServeTLS(listenAddress, tlsCert, tlsKey, nil))
	}

	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
