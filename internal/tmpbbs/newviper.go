package tmpbbs

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewViper() (*viper.Viper, error) {
	initFlags()
	pflag.Parse()

	viper := viper.New()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("tmpbbs")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, err
	}

	return viper, nil
}

func initFlags() {
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
	pflag.IntP("replies-per-page", "e", 10, "number of replies to show per page ($TMPBBS_REPLIES_PER_PAGE)") //nolint:mnd
	pflag.StringSliceP("css-urls", "u", []string{"/static/main.css"},
		"comma-separated list of CSS URLs ($TMPBBS_CSS_URLS)")
	pflag.BoolP("replies", "r", true, "enable replies ($TMPBBS_REPLIES)")
	pflag.BoolP("emoji", "m", true, "enable emoji shortcode expansion ($TMPBBS_EMOJI)")
	pflag.BoolP("qr-codes", "q", true, "enable shareable URL QR codes ($TMPBBS_QR_CODES)")
	pflag.BoolP("version", "v", false, "print version")
	pflag.BoolP("help", "h", false, "usage help")
}
