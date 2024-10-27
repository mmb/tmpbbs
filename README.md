Ephemeral internet forum software. All data is stored in memory only and when
the process ends all posts are gone.

# Features
  * Markdown support
  * optional TLS (using `-tls-cert` and `-tls-key` options)
  * trip code support in `username#secret` format

# Installation

```sh
go install github.com/mmb/tmpbbs@latest
```

# Options

```
% tmpbbs --help
Usage of tmpbbs:
  -l, --listen-address string   <host>:port to listen on ($TMPBBS_LISTEN_ADDRESS) (default ":8080")
  -c, --tls-cert string         path to PEM server certificate ($TMPBBS_TLS_CERT)
  -k, --tls-key string          path to PEM server key ($TMPBBS_TLS_KEY)
  -t, --title string            site title ($TMPBBS_TITLE) (default "tmpbbs")
  -a, --trip-code-salt string   random salt to use for generating trip codes ($TMPBBS_TRIP_CODE_SALT)
  -p, --load-posts string       path to YAML or JSON file of posts to load, format [{"title":"","author":"","body":""}] ($TMPBBS_LOAD_POSTS)
  -u, --css-urls strings        comma-separated list of CSS URLs ($TMPBBS_CSS_URLS) (default [/css])
  -h, --help                    usage help
```

# Compiling

```sh
go build
```
