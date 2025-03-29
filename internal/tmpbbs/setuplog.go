package tmpbbs

import (
	"log/slog"
	"os"
)

func SetupLog(jsonLog bool) {
	if jsonLog {
		logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
		slog.SetDefault(logger)
	}
}
