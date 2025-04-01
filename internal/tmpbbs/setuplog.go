package tmpbbs

import (
	"log/slog"
	"os"
)

// SetupLog configures the default slog logger.
func SetupLog(jsonLog bool) {
	if jsonLog {
		logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
		slog.SetDefault(logger)
	}
}
