package tmpbbs

import (
	"log/slog"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// LoggedViperSettings makes a map[string]any satisfy the slog LogValuer
// interface.
type LoggedViperSettings map[string]any

// NewViper returns a new viper.Viper with flags configured and the command
// line parsed.
func NewViper() (*viper.Viper, error) {
	initFlags()
	pflag.Parse()

	vipr := viper.New()
	vipr.SetEnvPrefix("tmpbbs")
	vipr.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	vipr.AutomaticEnv()

	err := vipr.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, err
	}

	configFile := vipr.GetString("config-file")
	if configFile != "" {
		vipr.SetConfigFile(configFile)

		err = vipr.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	return vipr, nil
}

// LogValue returns an [slog.GroupValue] with the subset of viper settings.
func (lvs LoggedViperSettings) LogValue() slog.Value {
	logConfigKeys := []string{
		"config-file",
		"css-urls",
		"emoji",
		"grpc-listen-address",
		"json-log",
		"listen-address",
		"load-posts",
		"prune-interval",
		"prune-max-age",
		"pull-interval",
		"pull-peers",
		"qr-codes",
		"replies-per-page",
		"replies",
		"serve-binary",
		"serve-fs-paths",
		"title",
		"tls-cert",
		"tls-key",
	}

	attrs := make([]slog.Attr, len(logConfigKeys))
	for i, key := range logConfigKeys {
		attrs[i] = slog.Any(key, lvs[key])
	}

	return slog.GroupValue(attrs...)
}

func initFlags() {
	pflag.StringP("config-file", "o", "", "path to config file ($TMPBBS_CONFIG_FILE)")
	pflag.StringSliceP("css-urls", "u", []string{}, "comma-separated list of additional CSS URLs ($TMPBBS_CSS_URLS)")
	pflag.BoolP("emoji", "m", true, "enable emoji shortcode expansion ($TMPBBS_EMOJI)")
	pflag.StringP("grpc-listen-address", "g", ":8081", "<host>:port to listen on for gRPC ($TMPBBS_GRPC_LISTEN_ADDRESS)")
	pflag.BoolP("json-log", "j", false, "set log output format to JSON ($TMPBBS_JSON_LOG)")
	pflag.StringP("listen-address", "l", ":8080", "<host>:port to listen on for HTTP ($TMPBBS_LISTEN_ADDRESS)")
	pflag.StringSliceP("load-posts", "p", []string{},
		`comma-separated paths of YAML or JSON files of posts to load, format [{"title":"","author":"","body":""}]
($TMPBBS_LOAD_POSTS)`)
	pflag.DurationP("prune-interval", "s", time.Hour,
		"how often to check for stale posts to prune ($TMPBBS_PRUNE_INTERVAL)")
	pflag.DurationP("prune-max-age", "w", 30*24*time.Hour,
		"delete posts that haven't been updated in this long ($TMPBBS_PRUNE_MAX_AGE)")
	pflag.DurationP("pull-interval", "i", 30*time.Second, //nolint:mnd // default
		"peer pull interval ($TMPBBS_PULL_INTERVAL)")
	pflag.StringSliceP("pull-peers", "d", []string{},
		"comma-separated list of tmpbbs gRPC <tls://>host:gRPCport to pull posts from ($TMPBBS_PULL_PEERS)")
	pflag.BoolP("qr-codes", "q", true, "enable shareable URL QR codes ($TMPBBS_QR_CODES)")
	pflag.BoolP("replies", "r", true, "enable replies ($TMPBBS_REPLIES)")
	pflag.IntP("replies-per-page", "e", 10, //nolint:mnd // default
		"number of replies to show per page ($TMPBBS_REPLIES_PER_PAGE)")
	pflag.BoolP("serve-binary", "b", true, "serve tmpbbs binary at /self ($TMPBBS_SERVE_BINARY)")
	pflag.StringSliceP("serve-fs-paths", "f", []string{},
		"comma-separated list of urlprefix=/local/dir to serve ($TMPBBS_SERVE_FS_PATHS)")
	pflag.StringSliceP("superuser-tripcodes", "n", []string{},
		"comma-separated list of superuser tripcodes ($TMPBBS_SUPERUSER_TRIPCODES)")
	pflag.StringP("title", "t", "tmpbbs", "site title ($TMPBBS_TITLE)")
	pflag.StringP("tls-cert", "c", "", "path to PEM server certificate ($TMPBBS_TLS_CERT)")
	pflag.StringP("tls-key", "k", "", "path to PEM server key ($TMPBBS_TLS_KEY)")
	pflag.StringP("tripcode-salt", "a", "", "random salt to use for generating tripcodes ($TMPBBS_TRIPCODE_SALT)")
	pflag.BoolP("version", "v", false, "print version")
}
