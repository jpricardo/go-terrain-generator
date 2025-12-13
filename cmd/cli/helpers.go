package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

const (
	outputDir = "./cmd/output/"
)

func saveBmp(img *image.Gray, fileName string) {
	f, err := os.Create(fmt.Sprintf("%s%s", outputDir, fileName))
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Panic(err)
	}
}
