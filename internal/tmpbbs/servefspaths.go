package tmpbbs

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// ErrInvalidServeFSPathsMapping is returned when a serve-fs-paths mapping is
// not in the format urlprefix=/local/dir.
var ErrInvalidServeFSPathsMapping = errors.New("serve-fs-paths mappings must be in the format urlprefix=/local/dir")

// ServeFSPaths adds an [http.Handler] to the passed in [http.ServeMux] that
// serves directories from the filesystem.
func ServeFSPaths(paths []string, serveMux *http.ServeMux) error {
	for _, dirMapping := range paths {
		urlPrefix, dir, found := strings.Cut(dirMapping, "=")
		if !found || urlPrefix == "" || dir == "" {
			return fmt.Errorf("%w: %q", ErrInvalidServeFSPathsMapping, dirMapping)
		}

		_, err := os.Stat(dir)
		if err != nil {
			return err
		}

		if !strings.HasPrefix(urlPrefix, "/") {
			urlPrefix = "/" + urlPrefix
		}

		serveMux.Handle(fmt.Sprintf("GET %s/", urlPrefix), http.StripPrefix(urlPrefix, http.FileServer(http.Dir(dir))))
	}

	return nil
}
