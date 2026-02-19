//nolint:godoclint // usage preformatted
/*
Tmpbbs serves a forum web site.

Forum posts are in-memory only and are lost if the process restarts.
Posts are anonymous but tripcodes are supported.
Markdown can be used in post bodies and emoji shortcodes can be used in all fields.
Posts can be bidirectionally synced from other tmpbbs instances over gRPC to create a distributed forum.
It can also serve directories of static files from the filesystem.
Posts can be loaded from a YAML file on startup.

Usage of tmpbbs:

	-o, --config-file string            path to config file ($TMPBBS_CONFIG_FILE)
	-u, --css-urls strings              comma-separated list of additional CSS URLs ($TMPBBS_CSS_URLS)
	-m, --emoji                         enable emoji shortcode expansion ($TMPBBS_EMOJI) (default true)
	-g, --grpc-listen-address string    <host>:port to listen on for gRPC ($TMPBBS_GRPC_LISTEN_ADDRESS) (default ":8081")
	-j, --json-log                      set log output format to JSON ($TMPBBS_JSON_LOG)
	-l, --listen-address string         <host>:port to listen on for HTTP ($TMPBBS_LISTEN_ADDRESS) (default ":8080")
	-p, --load-posts strings            comma-separated paths of YAML or JSON files of posts to load,
	                                    format [{"title":"","author":"","body":""}] ($TMPBBS_LOAD_POSTS)
	-s, --prune-interval duration       how often to check for stale posts to prune ($TMPBBS_PRUNE_INTERVAL)
	                                    (default 1h0m0s)
	-w, --prune-max-age duration        delete posts that haven't been updated in this long ($TMPBBS_PRUNE_MAX_AGE)
	                                    (default 720h0m0s)
	-i, --pull-interval duration        peer pull interval ($TMPBBS_PULL_INTERVAL) (default 30s)
	-d, --pull-peers strings            comma-separated list of tmpbbs gRPC <tls://>host:gRPCport to pull posts from
	                                    ($TMPBBS_PULL_PEERS)
	-q, --qr-codes                      enable shareable URL QR codes ($TMPBBS_QR_CODES) (default true)
	-r, --replies                       enable replies ($TMPBBS_REPLIES) (default true)
	-e, --replies-per-page int          number of replies to show per page ($TMPBBS_REPLIES_PER_PAGE) (default 10)
	-b, --serve-binary                  serve tmpbbs binary at /self ($TMPBBS_SERVE_BINARY) (default true)
	-f, --serve-fs-paths strings        comma-separated list of urlprefix=/local/dir to serve ($TMPBBS_SERVE_FS_PATHS)
	-n, --superuser-tripcodes strings   comma-separated list of superuser tripcodes ($TMPBBS_SUPERUSER_TRIPCODES)
	-t, --title string                  site title ($TMPBBS_TITLE) (default "tmpbbs")
	-c, --tls-cert string               path to PEM server certificate ($TMPBBS_TLS_CERT)
	-k, --tls-key string                path to PEM server key ($TMPBBS_TLS_KEY)
	-a, --tripcode-salt string          random salt to use for generating tripcodes ($TMPBBS_TRIPCODE_SALT)
	-v, --version                       print version
*/
package main

import (
	"context"
	"crypto/rand"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/mmb/tmpbbs/internal/tmpbbs"
)

// If commit is not set by go build (during development, for example) set it to the current time because it's used for
// caching.
var (
	commit  = strconv.FormatInt(time.Now().UnixNano(), 10) //nolint:gochecknoglobals // used with go build -X
	version = ""
	date    = "" //nolint:gochecknoglobals // used with go build -X
)

//go:embed static
var staticFS embed.FS

func main() {
	ctx := context.Background()

	viper, err := tmpbbs.NewViper()
	if err != nil {
		log.Fatal(err)
	}

	if viper.GetBool("version") {
		fmt.Fprintf(os.Stdout, "%s-%s-%s\n", version, commit, date)
		os.Exit(0)
	}

	tmpbbs.SetupLog(viper.GetBool("json-log"))

	slog.InfoContext(ctx, "startup", "version", version, "commit", commit, "date", date, "config",
		tmpbbs.LoggedViperSettings(viper.AllSettings()))

	postStore := tmpbbs.NewPostStore(viper.GetString("title"), viper.GetDuration("prune-interval"),
		viper.GetDuration("prune-max-age"))

	tripcoder, err := tmpbbs.NewTripcoder(viper.GetString("tripcode-salt"), viper.GetStringSlice("superuser-tripcodes"),
		rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	for _, loadPath := range viper.GetStringSlice("load-posts") {
		err = postStore.LoadYAML(loadPath, tripcoder)
		if err != nil {
			log.Fatal(err)
		}
	}

	if viper.GetString("grpc-listen-address") != "" {
		go func() {
			log.Fatal(tmpbbs.ServeGRPC(viper.GetString("grpc-listen-address"), viper.GetString("tls-cert"),
				viper.GetString("tls-key"), tmpbbs.NewPostSyncServer(postStore)))
		}()
	}

	err = tmpbbs.RunPullPeers(viper.GetStringSlice("pull-peers"), viper.GetDuration("pull-interval"), postStore)
	if err != nil {
		log.Fatal(err)
	}

	serveMux, err := tmpbbs.NewServeMux(viper, commit, staticFS, postStore, tripcoder)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(tmpbbs.Serve(viper.GetString("listen-address"), viper.GetString("tls-cert"), viper.GetString("tls-key"),
		serveMux))
}
