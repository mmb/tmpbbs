package tmpbbs

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// ServeFSPaths adds an [http.Handler] to the passed in [http.ServeMux] that
// serves directories from the filesystem.
func ServeFSPaths(paths []string, serveMux *http.ServeMux) error {
	for _, dirMapping := range paths {
		parts := strings.SplitN(dirMapping, "=", 2) //nolint:mnd // mapping has two parts, can't change
		urlPrefix, dir := parts[0], parts[1]

		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			return err
		}

		if !strings.HasPrefix(urlPrefix, "/") {
			urlPrefix = "/" + urlPrefix
		}

		serveMux.Handle(fmt.Sprintf("GET %s/", urlPrefix), http.StripPrefix(urlPrefix, http.FileServer(http.Dir(dir))))
	}

	return nil
}
