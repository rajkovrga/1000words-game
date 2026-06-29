package logger

import (
	"io"
	"log/slog"
	"os"
)

func New() *slog.Logger {
	_ = os.MkdirAll("storage/logs", 0755)

	file, err := os.OpenFile(
		"storage/logs/app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)

	var writer io.Writer = io.Discard

	if err == nil {
		writer = file
	}

	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}
