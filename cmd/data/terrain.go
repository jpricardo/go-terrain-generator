package data

import (
	"errors"
	"go-terrain-generator/cmd/helpers"
	"image"
	"image/color"
	"math"
)

type Terrain [][]*Point
type Chunk struct {
	X       int
	Y       int
	Terrain *Terrain
}

type Point struct {
	elevation uint8
	material  *Material
}

const (
	MinSize     = 32
	MaxSize     = 8192
	WaterLevel  = uint8(0)
	GroundLevel = uint8(96)
	SnowLevel   = uint8(192)
	UpperLimit  = uint8(255)
)

func NewChunk(x, y, size int) *Chunk {
	chunk := Chunk{
		X:       x,
		Y:       y,
		Terrain: NewTerrain(size),
	}

	return &chunk
}

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

	grassLevel := GroundLevel + 8
	stoneLevel := SnowLevel - (SnowLevel-GroundLevel)/3

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

			if e > stoneLevel {
				row[y].material = Stone()
				continue
			}

			if e < grassLevel {
				row[y].material = Sand()
				continue
			}

			row[y].material = Grass()
		}
	}

	return t, nil
}

func (t *Terrain) ApplyChunk(c *Chunk) *Terrain {
	startX := c.X
	startY := c.Y

	ct := c.Terrain

	for x := range *ct {
		for y := range *ct {
			(*t)[startX+x][startY+y] = (*ct)[x][y]
		}
	}

	return t
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
			c = color.RGBA(helpers.MergeColors(c, e, 0.5))

			img.Set(x, y, c)

		}
	}

	return img
}
