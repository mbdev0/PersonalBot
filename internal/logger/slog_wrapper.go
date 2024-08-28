package logger

import (
	"log/slog"
)

type Level = slog.Level

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

// this can be extended to include more types of values like int, float, etc. if needed
func String(key, value string) slog.Attr {
	return slog.Attr{Key: key, Value: slog.StringValue(value)}
}
