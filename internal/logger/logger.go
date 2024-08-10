package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	slogLogger *slog.Logger
)

func initSLog() {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	slogLogger = slog.New(jsonHandler)
}

func getLogger() *slog.Logger {
	if slogLogger == nil {
		initSLog()
	}
	return slogLogger
}

func Log(level slog.Level, msg string, attrs ...slog.Attr) {
	getLogger().LogAttrs(context.Background(), level, msg, attrs...)
}
