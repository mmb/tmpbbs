Ephemeral internet forum software. All data is stored in memory only and when
the process ends all posts are gone.

# Features
  * Markdown
  * emoji shortcodes
  * optional TLS (using `--tls-cert` and `--tls-key` options)
  * trip codes in `username#secret` format

# Installation

```sh
go install github.com/mmb/tmpbbs@latest
```

# Container Image

```sh
docker run --publish 8080:8080 mm62/tmpbbs
```

# Options

```
% tmpbbs --help
Usage of tmpbbs:
  -l, --listen-address string    <host>:port to listen on ($TMPBBS_LISTEN_ADDRESS) (default ":8080")
  -c, --tls-cert string          path to PEM server certificate ($TMPBBS_TLS_CERT)
  -k, --tls-key string           path to PEM server key ($TMPBBS_TLS_KEY)
  -t, --title string             site title ($TMPBBS_TITLE) (default "tmpbbs")
  -a, --trip-code-salt string    random salt to use for generating trip codes ($TMPBBS_TRIP_CODE_SALT)
  -p, --load-posts strings       comma-separated paths of YAML or JSON files of posts to load, format [{"title":"","author":"","body":""}] ($TMPBBS_LOAD_POSTS)
  -f, --serve-fs-paths strings   comma-separated list of urlprefix=/local/dir to serve ($TMPBBS_SERVE_FS_PATHS)
  -e, --replies-per-page int     Number of replies to show per page ($TMPBBS_REPLIES_PER_PAGE) (default 10)
  -u, --css-urls strings         comma-separated list of CSS URLs ($TMPBBS_CSS_URLS) (default [/static/main.css])
  -r, --replies                  Enable replies ($TMPBBS_REPLIES) (default true)
  -m, --emoji                    Enable emoji shortcode expansion ($TMPBBS_EMOJI) (default true)
  -h, --help                     usage help
```

# Compiling

```sh
go build
```
