package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func main() {
	inputDir := flag.String("i", "", "input directory")
	workers := flag.Int("c", runtime.NumCPU(), "number of workers")
	flag.Parse()

	if *inputDir == "" {
		fmt.Println("Input directory is required")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Convert flac to mp3 (inputDir=%s, workers=%d)\n", *inputDir, *workers)

	inputFiles, err := os.ReadDir(*inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read files from directory %s: %v\n", *inputDir, err)
		os.Exit(1)
	}

	outputDir := filepath.Join(filepath.Dir(*inputDir), filepath.Base(*inputDir)+"_320")
	if err := os.Mkdir(outputDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create output directory %s: %v\n", outputDir, err)
		os.Exit(1)
	}

	convertFiles(inputFiles, *inputDir, outputDir, *workers)
}

func convertFiles(inputFiles []os.DirEntry, inputDir, outputDir string, workers int) {
	defer track("convert all files")()
	bufferSize := 2 * workers
	inputFilesQueue := make(chan os.DirEntry, bufferSize)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(inputFilesQueue)
		for _, file := range inputFiles {
			if filepath.Ext(file.Name()) == ".flac" {
				inputFilesQueue <- file
			}
		}
	}()

	for range workers {
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
		ext := filepath.Ext(fileName)
		nameWithoutExt := fileName[:len(fileName)-len(ext)]

		fmt.Printf("Converting file: %s\n", nameWithoutExt)

		inputFilePath := filepath.Join(inputDir, fileName)
		outputFilePath := filepath.Join(outputDir, nameWithoutExt+".mp3")

		cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-ab", "320k", "-map_metadata", "0", "-id3v2_version", "3", outputFilePath)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to convert file %s: %v\n", inputFilePath, err)
		}
	}
}

func track(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("Execution time - %s: %v\n", name, time.Since(start))
	}
}
