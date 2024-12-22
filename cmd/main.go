package main

import (
	"log/slog"
	"os"
)

func init() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(logHandler))
}

func main() {
	slog.Info("Convert flac to mp3 ...")
}
