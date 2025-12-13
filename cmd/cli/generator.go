package main

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
	waterLevel   = uint8(64)
	maxElevation = uint8(255)
	lowResRatio  = float64(0.125)
)

func (t *Terrain) new(size int) Terrain {
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

func (t *Terrain) generateElevation(opts ElevationOptions) (Terrain, error) {
	if opts.smoothness > 1 || opts.smoothness < 0 {
		return nil, errors.New("invalid smoothness value")
	}

	var texture Texture
	et, _ := texture.whiteNoise(TextureOptions{size: len(*t), smoothness: opts.smoothness, baseLine: waterLevel})
	passes := int(math.Max(1/lowResRatio, 1))
	for p := 1; p <= passes; p++ {
		ls := float64(p) / lowResRatio
		nt, _ := texture.whiteNoise(TextureOptions{size: int(ls), smoothness: opts.smoothness, baseLine: waterLevel})
		et, _ = et.merge([]Texture{nt}, MergeOptions{smoothness: math.Pow(opts.smoothness, float64(p/passes)), opacity: math.Pow(.25, float64(passes/p))})

		img := et.toBitmap()
		saveBmp(img, fmt.Sprintf("elevation_texture_%d.png", p))
	}

	nt, _ := t.applyElevation(et, ApplyOptions{smoothness: opts.smoothness, opacity: opts.smoothness})
	return nt, nil
}

func generateTerrain(size int) (Terrain, error) {
	if size <= 0 || size > maxSize {
		return nil, errors.New("invalid map size")
	}

	fmt.Printf("[TERRAIN] Generating %dx%d terrain map...\n", size, size)

	var terrain Terrain
	terrain = terrain.new(size)
	terrain, err := terrain.generateElevation(ElevationOptions{smoothness: .5})
	if err != nil {
		return nil, err
	}

	return terrain, nil
}

func (t *Terrain) applyElevation(texture Texture, opts ApplyOptions) (Terrain, error) {
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

			row[y].elevation = uint8(math.Min(e, float64(maxElevation)))
		}
	}

	return *t, nil
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
