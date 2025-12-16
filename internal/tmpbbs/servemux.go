package tmpbbs

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/spf13/viper"
)

// NewServeMux returns a new http.ServeMux configured with routes and handlers
// based on user configuration.
func NewServeMux(vipr *viper.Viper, commit string, staticFS embed.FS, postStore *PostStore,
	tripcoder *Tripcoder,
) (http.Handler, error) {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/healthz", func(responseWriter http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(responseWriter, "ok")
	})

	postGetHandler := newPostGetHandler(postStore, commit, vipr.GetStringSlice("css-urls"), vipr.GetBool("emoji"),
		vipr.GetBool("replies"), vipr.GetInt("replies-per-page"), vipr.GetBool("qr-codes"))

	serveMux.Handle("GET /{$}", postGetHandler)
	serveMux.Handle("GET /p/{id}", postGetHandler)

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, err
	}

	serveMux.Handle("GET /static/", newImmutableGetHandler(http.StripPrefix("/static/"+commit,
		http.FileServerFS(staticDir))))
	serveMux.Handle("GET /robots.txt", http.FileServerFS(staticDir))

	if vipr.GetBool("replies") {
		postPostHandler := newPostPostHandler(postStore, tripcoder)
		serveMux.Handle("POST /{$}", postPostHandler)
		serveMux.Handle("POST /p/{parentID}", postPostHandler)
	}

	if vipr.GetBool("emoji") {
		serveMux.Handle("GET /emoji-suggest/", newImmutableGetHandler(newEmojiSuggestHandler()))
	}

	if vipr.GetBool("qr-codes") {
		serveMux.Handle("GET /qr/", newImmutableGetHandler(newQRCodeGetHandler()))
	}

	if vipr.GetBool("serve-binary") {
		serveMux.Handle("GET /self", newSelfGetHandler())
	}

	if pathsErr := ServeFSPaths(vipr.GetStringSlice("serve-fs-paths"), serveMux); pathsErr != nil {
		return nil, pathsErr
	}

	externalCSS := len(vipr.GetStringSlice("css-urls")) > 0

	return newCommonHeadersHandler(serveMux, externalCSS), nil
}
