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
)

var version = "unspecified"

//go:embed static
var staticFS embed.FS

func main() {
	tmpbbs.InitBluemonday()

	viper, err := tmpbbs.NewViper()
	if err != nil {
		log.Fatal(err)
	}

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
	repliesEnabled := viper.GetBool("replies")

	if repliesEnabled {
		postPostHandler := tmpbbs.NewPostPostHandler(repliesPerPage, postStore, tripCoder)
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

	log.Fatal(tmpbbs.Serve(viper.GetString("listen-address"), viper.GetString("tls-cert"), viper.GetString("tls-key")))
}
