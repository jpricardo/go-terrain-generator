package main

import (
	"fmt"
	"go-terrain-generator/cmd/data"
	"go-terrain-generator/cmd/helpers"
	"log"
)

const (
	outputDir      = "./cmd/output/"
	outputFileName = "output_grayscale.png"
)

func main() {
	size := 2048

	terrain, err := data.GenerateTerrain(size)
	if err != nil {
		log.Panic(err)
	}

	img := terrain.ToBitmap()
	helpers.SaveBmp(img, outputDir, outputFileName)

	fmt.Println("Terrain map generated succesfully!")
}
