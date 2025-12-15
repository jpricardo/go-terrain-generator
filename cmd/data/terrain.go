package data

import (
	"errors"
	"fmt"
	"go-terrain-generator/cmd/helpers"
	"image"
	"image/color"
	"math"
)

type Terrain [][]*Point

type Point struct {
	elevation uint8
	material  *Material
}

const (
	MinSize     = 32
	MaxSize     = 8192
	WaterLevel  = uint8(0)
	GroundLevel = uint8(64)
	SnowLevel   = uint8(96)
	UpperLimit  = uint8(255)
)

func NewTerrain(size int) *Terrain {
	var terrain Terrain

	for x := range size {
		terrain = append(terrain, make([]*Point, size))

		for y := range size {
			terrain[x][y] = &Point{elevation: GroundLevel - GroundLevel/4}
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

func (t *Terrain) ApplyMaterials() (*Terrain, error) {
	terrainSize := len(*t)

	for x := range terrainSize {
		row := (*t)[x]

		for y := range row {
			e := row[y].elevation

			if e <= GroundLevel {
				row[y].material = Water()
				continue
			}

			if e > SnowLevel {
				row[y].material = Snow()
				continue
			}

			if e > GroundLevel+(SnowLevel-GroundLevel)*3/4 {
				row[y].material = Stone()
				continue
			}

			if e < GroundLevel+4 {
				row[y].material = Sand()
				continue
			}

			row[y].material = Grass()
		}
	}

	return t, nil
}

func (t *Terrain) ToBitmap() *image.RGBA {
	size := len(*t)
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	for y := range size {
		for x := range size {
			point := (*t)[y][x]
			c := point.material.color
			a := point.elevation
			e := color.RGBA{
				R: a,
				G: a,
				B: a,
				A: 255,
			}
			c = color.RGBA(helpers.MergeColors(c, e))

			img.Set(x, y, c)

		}
	}

	return img
}
