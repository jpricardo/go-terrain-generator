package main

import (
	"fmt"
	"go-terrain-generator/cmd/helpers"
	"log"
	"time"
)

const (
	outputDir      = "./cmd/output/"
	outputFileName = "output.png"
)

func main() {
	size := 2048
	smoothness := 0.75
	seed := uint64(time.Now().Unix())

	terrain, err := GenerateTerrain(TerrainOptions{seed: seed, size: size, smoothness: smoothness})
	if err != nil {
		log.Panic(err)
	}

	img := terrain.ToBitmap()
	helpers.SaveBmp(img, outputDir, outputFileName)

	fmt.Println("Terrain map generated succesfully!")
}
