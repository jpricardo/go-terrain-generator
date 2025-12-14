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

	var texture Texture
	textures := []Texture{}
	et, _ := texture.WhiteNoise(TextureOptions{size: len(*t), smoothness: 1, baseLine: waterLevel})
	passes := int(math.Pow(2, 5))
	for p := range passes {
		ts := math.Max(float64(len(*t))/math.Pow(2, float64(passes-p)), 1)

		if ts < minSize {
			continue
		}

		nt, _ := texture.ValueNoise(TextureOptions{
			size:       int(ts),
			smoothness: math.Pow(opts.smoothness, 1),
			baseLine:   waterLevel,
		})
		nt = nt.ApplyFilters([]Filter{Contrast(float64(p) / float64(passes))})

		img := nt.ToBitmap()
		helpers.SaveBmp(img, outputDir, fmt.Sprintf("elevation_texture_%d.png", p))

		textures = append(textures, nt)
	}

	mt, _ := et.Merge(textures, MergeOptions{
		opacity:    0.1,
		smoothness: 1 - opts.smoothness,
	})
	mt = mt.ApplyFilters([]Filter{Contrast(1 - opts.smoothness)})
	nt, _ := t.ApplyElevation(mt, ApplyOptions{
		smoothness: opts.smoothness,
		opacity:    1 - math.Pow(opts.smoothness, 2),
	},
	)
	return nt, nil
}

func GenerateTerrain(size int) (Terrain, error) {
	if size <= 0 || size > maxSize || size < minSize {
		return nil, errors.New("invalid map size")
	}

	fmt.Printf("[TERRAIN] Generating %dx%d terrain map...\n", size, size)

	var terrain Terrain
	terrain = terrain.New(size)
	terrain, err := terrain.GenerateElevation(ElevationOptions{smoothness: .5})
	if err != nil {
		return nil, err
	}

	return terrain, nil
}

type ApplyOptions struct {
	smoothness float64
	opacity    float64
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
			tc := td * float64(te) * opts.smoothness * opts.opacity
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
