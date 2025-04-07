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

	outputDir := filepath.Join(filepath.Dir(*inputDir), filepath.Base(*inputDir)+"_320")

	if err := os.Mkdir(outputDir, 0o755); err != nil {
		slog.Error("Cannot create output directory", "outputDir", outputDir)
	}

	convertFiles(inputFiles, *inputDir, outputDir, *workers)
}

func convertFiles(inputFiles []os.DirEntry, inputDir, outputDir string, workers int) {
	bufferSize := 2 * workers
	inputFilesQueue := make(chan os.DirEntry, bufferSize)

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer close(inputFilesQueue)
		defer wg.Done()
		for _, file := range inputFiles {
			if filepath.Ext(file.Name()) == ".flac" {
				inputFilesQueue <- file
			}
		}
	}()

	for w := 0; w < workers; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			conversionWorker(inputFilesQueue, inputDir, outputDir)
		}()
	}

	wg.Wait()
}

func conversionWorker(queue chan os.DirEntry, inputDir, outputDir string) {
	for file := range queue {
		fileName := file.Name()
		fileExtension := filepath.Ext(fileName)
		fileNameWithoutExtension := fileName[0 : len(fileName)-len(fileExtension)]

		slog.Info("Convert file", "filename", fileNameWithoutExtension)

		inputFilePath := filepath.Join(inputDir, fileName)

		outputFilePath := filepath.Join(outputDir, fileNameWithoutExtension+".mp3")

		cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-ab", "320k", "-map_metadata", "0", "-id3v2_version", "3", outputFilePath)

		if err := cmd.Run(); err != nil {
			slog.Error("Failed to convert file", "inputFilePath", inputDir, "error", err)
		}
	}
}
