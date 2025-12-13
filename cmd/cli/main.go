package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
)

const (
	maxSize        = 2048
	outputDir      = "./cmd/output/"
	outputFileName = "output_grayscale.png"
)

func main() {
	size := 512

	terrain, err := generateTerrain(size)
	if err != nil {
		log.Panic(err)
	}

	img := terrain.toBitmap()

	f, err := os.Create(fmt.Sprintf("%s%s", outputDir, outputFileName))
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Panic(err)
	}

	fmt.Println("Terrain map generated succesfully!")
}
