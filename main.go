package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func convert(src, dst string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	fmt.Printf("Converting %s to %s\n", src, dst)
	pngFile, err := os.Open(src)
	if err != nil {
		fmt.Printf("Error Opening %v: %v\n", src, err)
		return
	}
	defer pngFile.Close()

	pngImage, err := png.Decode(pngFile)
	if err != nil {
		fmt.Printf("Error Decoding %v: %v\n", src, err)
		return
	}

	jpegFile, err := os.Create(dst)
	if err != nil {
		fmt.Printf("Error Opening %v: %v\n", dst, err)
		return
	}
	defer jpegFile.Close()

	err = jpeg.Encode(jpegFile, pngImage, nil)
	if err != nil {
		fmt.Printf("Error Encoding %v: %v\n", dst, err)
		return
	}
	if wg != nil {
		fmt.Printf("Converted %s to %s\n", src, dst)
	}
}

func main() {
	var (
		inputDir  string
		outputDir string
		parallel  bool
		h         bool
		help      bool
	)

	flag.StringVar(&inputDir, "i", "", "Input directory")
	flag.StringVar(&outputDir, "o", "", "Output directory")
	flag.BoolVar(&parallel, "p", false, "Run in parallel")
	flag.BoolVar(&h, "h", false, "Print help")
	flag.BoolVar(&help, "help", false, "Print help")
	flag.Parse()

	if h || help {
		flag.PrintDefaults()
		return
	}
	if inputDir == "" || outputDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("Could not read directory: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, file := range entries {
		if file.IsDir() {
			continue
		}
		input := filepath.Join(inputDir, file.Name())
		output := filepath.Join(outputDir, strings.TrimSuffix(file.Name(), ".png")+".jpg")
		if parallel {
			wg.Add(1)
			go convert(input, output, &wg)
		} else {
			convert(input, output, nil)
		}
	}
	if parallel {
		wg.Wait()
	}
}
