package main

import (
	"crypto/rand"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/mmb/tmpbbs/internal/tmpbbs"
	"github.com/spf13/pflag"
)

var version = "unspecified"

//go:embed static
var staticFS embed.FS

func main() {
	viper, err := tmpbbs.NewViper()
	if err != nil {
		log.Fatal(err)
	}

	if viper.GetBool("help") {
		pflag.CommandLine.SortFlags = false
		pflag.Usage()
		os.Exit(0)
	}

	if viper.GetBool("version") {
		fmt.Println(version)
		os.Exit(0)
	}

	postStore := tmpbbs.NewPostStore(viper.GetString("title"))

	tripcoder, err := tmpbbs.NewTripcoder(viper.GetString("tripcode-salt"), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	for _, loadPath := range viper.GetStringSlice("load-posts") {
		if err = postStore.LoadYAML(loadPath, tripcoder); err != nil {
			log.Fatal(err)
		}
	}

	if viper.GetString("grpc-listen-address") != "" {
		go func() {
			log.Fatal(tmpbbs.ServeGRPC(viper.GetString("grpc-listen-address"), viper.GetString("tls-cert"),
				viper.GetString("tls-key"), tmpbbs.NewPostSyncServer(postStore)))
		}()
	}

	serveMux, err := tmpbbs.NewServeMux(viper, staticFS, postStore, tripcoder)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(tmpbbs.Serve(viper.GetString("listen-address"), viper.GetString("tls-cert"), viper.GetString("tls-key"),
		serveMux))
}
