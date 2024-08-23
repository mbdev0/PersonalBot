package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
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
	pc, file, line, ok := runtime.Caller(1)
	if ok && level == slog.LevelError {
		attrs[len(attrs)-1] = slog.String("Stack:", fmt.Sprintf("%s File:%s:%d", runtime.FuncForPC(pc).Name(), file, line))
	}
	getLogger().LogAttrs(context.Background(), level, msg, attrs...)
}
