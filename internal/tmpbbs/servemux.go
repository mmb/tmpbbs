package tmpbbs

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/spf13/viper"
)

// NewServeMux returns a new http.ServeMux configured with routes and handlers
// based on user configuration.
func NewServeMux(vipr *viper.Viper, staticFS embed.FS, postStore *PostStore,
	tripcoder *Tripcoder,
) (*http.ServeMux, error) {
	serveMux := http.NewServeMux()

	postGetHandler := newPostGetHandler(vipr.GetInt("replies-per-page"), vipr.GetStringSlice("css-urls"),
		vipr.GetBool("replies"), vipr.GetBool("emoji"), vipr.GetBool("qr-codes"), postStore)
	serveMux.Handle("GET /{$}", postGetHandler)
	serveMux.Handle("GET /p/{uuid}", postGetHandler)

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, err
	}

	serveMux.Handle("GET /static/", newImmutableGetHandler(http.StripPrefix("/static", http.FileServerFS(staticDir))))
	serveMux.Handle("GET /robots.txt", newImmutableGetHandler(http.FileServerFS(staticDir)))

	if vipr.GetBool("replies") {
		postPostHandler := newPostPostHandler(postStore, tripcoder)
		serveMux.Handle("POST /{$}", postPostHandler)
		serveMux.Handle("POST /p/{parentUUID}", postPostHandler)
	}

	if vipr.GetBool("emoji") {
		serveMux.Handle("GET /emoji-suggest", newImmutableGetHandler(newEmojiSuggestHandler()))
	}

	if vipr.GetBool("qr-codes") {
		serveMux.Handle("GET /qr", newImmutableGetHandler(newQRCodeGetHandler()))
	}

	if vipr.GetBool("serve-binary") {
		serveMux.Handle("GET /self", newImmutableGetHandler(newSelfGetHandler()))
	}

	if pathsErr := ServeFSPaths(vipr.GetStringSlice("serve-fs-paths"), serveMux); pathsErr != nil {
		return nil, pathsErr
	}

	return serveMux, nil
}
