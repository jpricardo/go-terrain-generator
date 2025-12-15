package data

import (
	"errors"
	"fmt"
	"image"
	"math"
)

type Terrain [][]*Point

type Point struct {
	elevation uint8
}

const (
	MinSize    = 32
	MaxSize    = 8192
	WaterLevel = uint8(64)
)

func NewTerrain(size int) *Terrain {
	var terrain Terrain

	for x := range size {
		terrain = append(terrain, make([]*Point, size))

		for y := range size {
			terrain[x][y] = &Point{elevation: WaterLevel}
		}
	}

	return &terrain
}

type ApplyOptions struct {
	Opacity float64
}

func (t *Terrain) ApplyElevation(texture *Texture, opts ApplyOptions) (*Terrain, error) {
	terrainSize := len(*t)
	textureSize := len(*texture)

	if textureSize <= 0 {
		return nil, errors.New("invalid texture size")
	}

	fmt.Printf("[TEXTURE] Applying %dx%d elevation texture to %dx%d terrain...\n", textureSize, textureSize, terrainSize, terrainSize)

	p := terrainSize / textureSize
	for x := range terrainSize {
		row := (*t)[x]

		for y := range row {
			c := row[y].elevation

			lx := int(math.Min(float64(x/p), float64(textureSize-1)))
			ly := int(math.Min(float64(y/p), float64(textureSize-1)))

			te := (*texture)[lx][ly]
			td := 1 - float64(c)/float64(te)
			tc := td * float64(te) * opts.Opacity
			e := float64(c) + tc

			row[y].elevation = uint8(e)
		}
	}

	return t, nil
}

func (t *Terrain) ToBitmap() *image.Gray {
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
