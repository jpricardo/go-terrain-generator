package main

import (
	"fmt"
	"log"
)

const (
	maxSize        = 2048
	outputFileName = "output_grayscale.png"
)

func main() {
	size := 2048

	terrain, err := generateTerrain(size)
	if err != nil {
		log.Panic(err)
	}

	img := terrain.toBitmap()
	saveBmp(img, outputFileName)

	fmt.Println("Terrain map generated succesfully!")
}
