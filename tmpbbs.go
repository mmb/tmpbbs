package main

import (
	"crypto/rand"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/mmb/tmpbbs/internal/tmpbbs"
	"github.com/spf13/pflag"
)

var version = "unspecified"

//go:embed static
var staticFS embed.FS

func main() {
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

	tripcoder, err := tmpbbs.NewTripcoder(viper.GetString("tripcode-salt"), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	for _, loadPath := range viper.GetStringSlice("load-posts") {
		if err = postStore.LoadYAML(loadPath, tripcoder); err != nil {
			log.Fatal(err)
		}
	}

	repliesPerPage := viper.GetInt("replies-per-page")
	if viper.GetBool("replies") {
		postPostHandler := tmpbbs.NewPostPostHandler(repliesPerPage, postStore, tripcoder)
		http.Handle("POST /{$}", postPostHandler)
		http.Handle("POST /{parentID}", postPostHandler)
	}

	if viper.GetBool("qr-codes") {
		http.Handle("GET /qr", tmpbbs.NewQRCodeGetHandler())
	}

	postGetHandler := tmpbbs.NewPostGetHandler(repliesPerPage, viper.GetStringSlice("css-urls"), viper.GetBool("replies"),
		viper.GetBool("emoji"), viper.GetBool("qr-codes"), postStore)
	http.Handle("GET /{$}", postGetHandler)
	http.Handle("GET /{id}", postGetHandler)

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("GET /static/", http.StripPrefix("/static", http.FileServerFS(staticDir)))
	http.Handle("GET /robots.txt", http.FileServerFS(staticDir))

	if err = tmpbbs.ServeFSPaths(viper.GetStringSlice("serve-fs-paths")); err != nil {
		log.Fatal(err)
	}

	log.Fatal(tmpbbs.Serve(viper.GetString("listen-address"), viper.GetString("tls-cert"), viper.GetString("tls-key")))
}
