tmpbbs is an ephemeral, anonymous forum site. All posts are stored in memory
and lost when the process stops. It's a single static binary and does not require the
internet or any disk. It runs on a wide variety of hardware including
older wireless routers.

# Features
  * Markdown
  * emoji shortcodes
  * optional TLS (using `--tls-cert` and `--tls-key` options)
  * tripcodes in `username#secret` format
  * shareable URL QR codes

# Installation

```sh
go install github.com/mmb/tmpbbs@latest
```

# Container Images

```sh
docker run --publish 8080:8080 mm62/tmpbbs
```

or

```sh
docker run --publish 8080:8080 ghcr.io/mmb/tmpbbs
```

# Options

```
% tmpbbs --help
Usage of tmpbbs:
  -l, --listen-address string        <host>:port to listen on for HTTP ($TMPBBS_LISTEN_ADDRESS) (default ":8080")
  -g, --grpc-listen-address string   <host>:port to listen on for gRPC ($TMPBBS_GRPC_LISTEN_ADDRESS) (default ":8081")
  -c, --tls-cert string              path to PEM server certificate ($TMPBBS_TLS_CERT)
  -k, --tls-key string               path to PEM server key ($TMPBBS_TLS_KEY)
  -t, --title string                 site title ($TMPBBS_TITLE) (default "tmpbbs")
  -a, --tripcode-salt string         random salt to use for generating tripcodes ($TMPBBS_TRIPCODE_SALT)
  -p, --load-posts strings           comma-separated paths of YAML or JSON files of posts to load, format [{"title":"","author":"","body":""}]
                                     ($TMPBBS_LOAD_POSTS)
  -f, --serve-fs-paths strings       comma-separated list of urlprefix=/local/dir to serve ($TMPBBS_SERVE_FS_PATHS)
  -e, --replies-per-page int         number of replies to show per page ($TMPBBS_REPLIES_PER_PAGE) (default 10)
  -u, --css-urls strings             comma-separated list of CSS URLs ($TMPBBS_CSS_URLS) (default [/static/main.css])
  -r, --replies                      enable replies ($TMPBBS_REPLIES) (default true)
  -m, --emoji                        enable emoji shortcode expansion ($TMPBBS_EMOJI) (default true)
  -q, --qr-codes                     enable shareable URL QR codes ($TMPBBS_QR_CODES) (default true)
  -b, --serve-binary                 serve tmpbbs binary at /self ($TMPBBS_SERVE_BINARY) (default true)
  -v, --version                      print version
  -h, --help                         usage help
```

# Compiling

```sh
go build
```
