package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sync"
)

var (
	slogLogger *slog.Logger
	once       sync.Once
)

const (
	stackSkip = 2
)

func initSLog() {
	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		slogLogger = slog.New(jsonHandler)
		Log(slog.LevelError, "Error opening log file", ErrorMessage(err))
	}
	jsonHandler := slog.NewJSONHandler(io.MultiWriter(os.Stderr, f), nil)
	slogLogger = slog.New(jsonHandler)
}

func getLogger() *slog.Logger {
	once.Do(func() {
		initSLog()
		slog.SetDefault(slogLogger)
	})

	return slogLogger
}
func logInternal(level slog.Level, msg string, attrs ...slog.Attr) {
	_, file, line, ok := runtime.Caller(stackSkip)
	if ok && level == slog.LevelError {
		var lastElement int
		if len(attrs) == 0 {
			lastElement = 0
			attrs = make([]slog.Attr, 1)
		} else {
			attrs = append(attrs, slog.Attr{})
			lastElement = len(attrs) - 1
		}
		attrs[lastElement] = StringMessage("stack:", fmt.Sprintf("%s:%d", file, line))
	}
	getLogger().LogAttrs(context.Background(), level, msg, attrs...)
}

func Information(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logInternal(LevelInfo, msg)
}

func Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logInternal(LevelError, msg)
}

func Log(level slog.Level, msg string, attrs ...slog.Attr) {
	logInternal(level, msg, attrs...)
}
