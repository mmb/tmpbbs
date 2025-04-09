package tmpbbs

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/spf13/viper"
)

// NewServeMux returns a new http.ServeMux configured with routes and handlers
// based on user configuration.
func NewServeMux(viper *viper.Viper, staticFS embed.FS, postStore *PostStore,
	tripcoder *Tripcoder,
) (*http.ServeMux, error) {
	serveMux := http.NewServeMux()

	postGetHandler := newPostGetHandler(viper.GetInt("replies-per-page"), viper.GetStringSlice("css-urls"),
		viper.GetBool("replies"), viper.GetBool("emoji"), viper.GetBool("qr-codes"), postStore)
	serveMux.Handle("GET /{$}", postGetHandler)
	serveMux.Handle("GET /p/{uuid}", postGetHandler)

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, err
	}

	serveMux.Handle("GET /static/", newImmutableGetHandler(http.StripPrefix("/static", http.FileServerFS(staticDir)),
		postStore.posts[0].time))
	serveMux.Handle("GET /robots.txt", newImmutableGetHandler(http.FileServerFS(staticDir), postStore.posts[0].time))

	if viper.GetBool("replies") {
		postPostHandler := newPostPostHandler(viper.GetInt("replies-per-page"), postStore, tripcoder)
		serveMux.Handle("POST /{$}", postPostHandler)
		serveMux.Handle("POST /p/{parentUUID}", postPostHandler)
	}

	if viper.GetBool("emoji") {
		serveMux.Handle("GET /emoji-suggest", newImmutableGetHandler(newEmojiSuggestHandler(), postStore.posts[0].time))
	}

	if viper.GetBool("qr-codes") {
		serveMux.Handle("GET /qr", newImmutableGetHandler(newQRCodeGetHandler(), postStore.posts[0].time))
	}

	if viper.GetBool("serve-binary") {
		serveMux.Handle("GET /self", newImmutableGetHandler(newSelfGetHandler(), postStore.posts[0].time))
	}

	if pathsErr := ServeFSPaths(viper.GetStringSlice("serve-fs-paths"), serveMux); pathsErr != nil {
		return nil, pathsErr
	}

	return serveMux, nil
}
