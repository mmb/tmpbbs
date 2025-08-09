package tmpbbs

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

// SetupLog configures the default slog logger.
func SetupLog(jsonLog bool) {
	var handler slog.Handler

	if jsonLog {
		handler = slog.NewJSONHandler(os.Stderr, nil)
	} else {
		handler = log.New(os.Stderr)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
