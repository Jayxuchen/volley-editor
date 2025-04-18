package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jayxuchen/volley-editor/video"
)

func main() {
	inputDir := "input"
	outputDir := "output"

	// Create output directory if needed
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}

	// Scan input directory for video files
	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("failed to read input dir: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if strings.ToLower(ext) != ".mp4" && strings.ToLower(ext) != ".mov" {
			fmt.Println("ignoring file" + file.Name())
			continue
		}

		inputPath := filepath.Join(inputDir, file.Name())
		nameOnly := strings.TrimSuffix(file.Name(), ext)
		outputPath := filepath.Join(outputDir, nameOnly+"_edit"+ext)

		// Example overlays - you can replace with dynamic or loaded-from-file data
		overlays := []video.ScoreOverlay{
			{Start: 0, End: 1, Text: "Home 0 - Away 0"},
			{Start: 1, End: 2, Text: "Home 1 - Away 0"},
			{Start: 2, End: 3, Text: "Home 1 - Away 1"},
			{Start: 2, End: 3, Text: "Home 1 - Away 2"},
		}

		log.Printf("Processing %s → %s", inputPath, outputPath)
		if err := video.ApplyScoreOverlay(inputPath, outputPath, overlays); err != nil {
			log.Printf("Error processing %s: %v", inputPath, err)
		} else {
			log.Printf("✅ Wrote %s", outputPath)
		}
	}
}
