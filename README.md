Ephemeral internet forum software. All data is stored in memory only and when
the process ends all posts are gone.

# Features
  * optional TLS using (`-tlsCert` and `-tlsKey` options)
  * trip code support in `username#secret` format if `-tripCodeSalt` is provided

# Installation

```sh
go install github.com/mmb/tmpbbs@latest
```

# Options

```
% ./tmpbbs --help
Usage of ./tmpbbs:
  -cssURL string
        CSS URL (default "/css")
  -listenAddress string
        <host>:port to listen on (default ":8080")
  -title string
        Site title (default "tmpbbs")
  -tlsCert string
        path to PEM server certificate
  -tlsKey string
        path to PEM server key
  -tripCodeSalt string
        random salt to use for generating trip codes
```

# Compiling

```sh
go build
```
