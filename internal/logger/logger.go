package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
)

var (
	slogLogger *slog.Logger
)

func initSLog() {
	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		slogLogger = slog.New(jsonHandler)
		Log(slog.LevelError, "Error opening log file", Error(err))
	}
	jsonHandler := slog.NewJSONHandler(io.MultiWriter(os.Stderr, f), nil)
	slogLogger = slog.New(jsonHandler)
}

func getLogger() *slog.Logger {
	if slogLogger == nil {
		initSLog()
	}
	slog.SetDefault(slogLogger)
	return slogLogger
}

func Log(level slog.Level, msg string, attrs ...slog.Attr) {
	pc, file, line, ok := runtime.Caller(1)
	if ok && level == slog.LevelError {

		var lastElement int
		if len(attrs) == 0 {
			lastElement = 0
			attrs = make([]slog.Attr, 1)
		} else {
			attrs = append(attrs, slog.Attr{})
			lastElement = len(attrs) - 1
		}
		attrs[lastElement] = String("stack:", fmt.Sprintf("%s File:%s:%d", runtime.FuncForPC(pc).Name(), file, line))
	}

	getLogger().LogAttrs(context.Background(), level, msg, attrs...)
}
