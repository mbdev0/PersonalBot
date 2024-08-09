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

func GetLogger() *slog.Logger {
	if slogLogger == nil {
		initSLog()
	}
	return slogLogger
}

func Log(level slog.Level, msg string, attrs ...slog.Attr) {
	GetLogger().LogAttrs(context.Background(), level, msg, attrs...)
}
