package data

import (
	"errors"
	"fmt"
	"go-terrain-generator/cmd/helpers"
	"image"
	"math"
)

type Terrain [][]*Point

type Point struct {
	elevation uint8
}

const (
	minSize    = 32
	maxSize    = 8192
	waterLevel = uint8(64)
	outputDir  = "./cmd/output/"
)

func (t *Terrain) New(size int) Terrain {
	var terrain Terrain

	for x := range size {
		terrain = append(terrain, make([]*Point, size))

		for y := range size {
			terrain[x][y] = &Point{elevation: waterLevel}
		}
	}

	return terrain
}

type ElevationOptions struct {
	smoothness float64
}

func (t *Terrain) GenerateElevation(opts ElevationOptions) (Terrain, error) {
	if opts.smoothness > 1 || opts.smoothness < 0 {
		return nil, errors.New("invalid smoothness value")
	}

	passes := 4
	scale := (1 - opts.smoothness) * minSize / float64(len(*t))
	nt := PerlinNoise(PerlinNoiseOptions{
		size:   len(*t),
		scale:  scale,
		passes: passes,
	})

	img := nt.ToBitmap()
	helpers.SaveBmp(img, outputDir, "elevation_texture.png")

	return t.ApplyElevation(nt, ApplyOptions{opacity: 0.5 - (0.25 * opts.smoothness)})
}

func GenerateTerrain(size int) (Terrain, error) {
	if size <= 0 || size > maxSize || size < minSize {
		return nil, errors.New("invalid map size")
	}

	fmt.Printf("[TERRAIN] Generating %dx%d terrain map...\n", size, size)

	var terrain Terrain
	terrain = terrain.New(size)
	terrain, err := terrain.GenerateElevation(ElevationOptions{smoothness: 0.25})
	if err != nil {
		return nil, err
	}

	return terrain, nil
}

type ApplyOptions struct {
	opacity float64
}

func (t *Terrain) ApplyElevation(texture Texture, opts ApplyOptions) (Terrain, error) {
	terrainSize := len(*t)
	textureSize := len(texture)

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

			te := texture[lx][ly]
			td := 1 - float64(c)/float64(te)
			tc := td * float64(te) * opts.opacity
			e := float64(c) + tc

			row[y].elevation = uint8(e)
		}
	}

	return *t, nil
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
