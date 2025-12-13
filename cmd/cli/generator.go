package main

import (
	"errors"
	"fmt"
	"image"
)

type Terrain [][]*Point

type Point struct {
	elevation uint8
}

const (
	waterLevel = uint8(64)
)

func initializeTerrain(size int) Terrain {
	var terrain Terrain

	for x := range size {
		terrain = append(terrain, []*Point{})

		for range size {
			terrain[x] = append(terrain[x], &Point{})
		}
	}

	return terrain
}

func (t *Terrain) generateElevation() Terrain {
	for x := range *t {
		row := (*t)[x]

		for y := range row {
			point := row[y]
			point.elevation = waterLevel
		}
	}

	return *t
}

func generateTerrain(size int) (Terrain, error) {
	if size <= 0 || size > maxSize {
		return nil, errors.New("invalid map size")
	}

	fmt.Printf("Generating %dx%d terrain map...\n", size, size)

	terrain := initializeTerrain(size)
	terrain = terrain.generateElevation()

	return terrain, nil
}

func (t *Terrain) toBitmap() *image.Gray {
	size := len(*t)
	img := image.NewGray(image.Rect(0, 0, size, size))

	for y := range size {
		for x := range size {
			offset := (y * img.Stride) + x
			img.Pix[offset] = (*t)[y][x].elevation
		}
	}

	return img
}

func (t *Terrain) print() {
	for x := range *t {
		row := (*t)[x]

		for y := range row {
			point := row[y]
			fmt.Print(point)
		}

		fmt.Print('\n')
	}
}
