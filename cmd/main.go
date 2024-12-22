package main

import (
	"flag"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

func init() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(logHandler))
}

func main() {
	inputDir := (flag.String("i", "", "input directory"))
	workers := flag.Int("c", runtime.NumCPU(), "number of workers")

	flag.Parse()

	slog.Info("Convert flac to mp3", "inputDir", inputDir, "workers", workers)

	inputFiles, err := os.ReadDir(*inputDir)

	if err != nil {
		slog.Error("Cannot read files", "inputDir", inputDir, "error", err)
	}

	semaphore := make(chan struct{}, *workers)
	var wg sync.WaitGroup

	for _, file := range inputFiles {
		fileName := file.Name()
		fileExtension := filepath.Ext(fileName)
		if fileExtension == ".flac" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				fileNameWithoutExtension := fileName[0 : len(fileName)-len(fileExtension)]

				slog.Info("Convert file", "filename", fileNameWithoutExtension)

				inputFilePath := filepath.Join(*inputDir, fileName)

				outputDir := filepath.Join(filepath.Dir(*inputDir), filepath.Base(*inputDir)+"_320")

				if err := os.Mkdir(outputDir, 0755); err != nil {
					slog.Error("Cannot create output directory", "outputDir", outputDir)
				}

				outputFilePath := filepath.Join(outputDir, fileNameWithoutExtension+".mp3")

				cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-ab", "320k", "-map_metadata", "0", "-id3v2_version", "3", outputFilePath)

				if err := cmd.Run(); err != nil {
					slog.Error("Failed to convert file", "inputFilePath", inputDir, "error", err)
				}
			}()
		}
	}
	wg.Wait()
}
