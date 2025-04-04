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

	postGetHandler := NewPostGetHandler(viper.GetInt("replies-per-page"), viper.GetStringSlice("css-urls"),
		viper.GetBool("replies"), viper.GetBool("emoji"), viper.GetBool("qr-codes"), postStore)
	serveMux.Handle("GET /{$}", postGetHandler)
	serveMux.Handle("GET /p/{uuid}", postGetHandler)

	staticDir, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, err
	}

	serveMux.Handle("GET /static/", http.StripPrefix("/static", http.FileServerFS(staticDir)))
	serveMux.Handle("GET /robots.txt", http.FileServerFS(staticDir))

	if viper.GetBool("replies") {
		postPostHandler := NewPostPostHandler(viper.GetInt("replies-per-page"), postStore, tripcoder)
		serveMux.Handle("POST /{$}", postPostHandler)
		serveMux.Handle("POST /p/{parentUUID}", postPostHandler)
	}

	if viper.GetBool("qr-codes") {
		serveMux.Handle("GET /qr", NewQRCodeGetHandler())
	}

	if viper.GetBool("serve-binary") {
		serveMux.Handle("GET /self", NewSelfGetHandler())
	}

	if pathsErr := ServeFSPaths(viper.GetStringSlice("serve-fs-paths"), serveMux); pathsErr != nil {
		return nil, pathsErr
	}

	return serveMux, nil
}
