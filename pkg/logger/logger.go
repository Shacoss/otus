package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once
)

func GetLogger() *slog.Logger {
	once.Do(func() {
		// Используем JSON-обработчик
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		instance = slog.New(handler)
	})
	return instance
}
