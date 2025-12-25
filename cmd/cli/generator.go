package main

import (
	"errors"
	"go-terrain-generator/cmd/data"
	"go-terrain-generator/cmd/helpers"
	"image"
	"image/color"
	"log"
)

type ChunkOptions struct {
	x          int
	y          int
	seed       uint64
	size       int
	smoothness float64
}

func GenerateChunk(img *image.RGBA, ch chan *data.Chunk, opts ChunkOptions) {
	cx := opts.x * opts.size
	cy := opts.y * opts.size

	c := data.NewChunk(cx, cy, opts.size)
	ct := c.Terrain

	et, err := GenerateElevation(ElevationOptions{
		x:          cx,
		y:          cy,
		seed:       opts.seed,
		smoothness: opts.smoothness,
		size:       opts.size,
	})
	if err != nil {
		log.Panic(err)
	}

	ct, err = ct.ApplyElevation(et, data.ApplyOptions{Opacity: 1})
	if err != nil {
		log.Panic(err)
	}

	ct, err = ct.ApplyMaterials()
	if err != nil {
		log.Panic(err)
	}

	WriteBmp(img, c)

	ch <- c
}

type TerrainOptions struct {
	seed       uint64
	size       int
	smoothness float64
}

func GenerateTerrain(opts TerrainOptions) (*image.RGBA, error) {
	if opts.size <= 0 || opts.size > data.MaxSize || opts.size < data.MinSize {
		return nil, errors.New("invalid map size")
	}

	chunkSize := data.MinSize
	chunks := opts.size / chunkSize
	ch := make(chan *data.Chunk)
	img := image.NewRGBA(image.Rect(0, 0, opts.size, opts.size))

	for x := range chunks {
		for y := range chunks {
			go GenerateChunk(img, ch, ChunkOptions{
				x:          x,
				y:          y,
				seed:       opts.seed,
				smoothness: opts.smoothness,
				size:       chunkSize,
			})
		}
	}

	for range chunks * chunks {
		<-ch
	}

	return img, nil
}

type ElevationOptions struct {
	seed       uint64
	x          int
	y          int
	size       int
	smoothness float64
}

func GenerateElevation(opts ElevationOptions) (*data.Texture, error) {
	if opts.smoothness > 1 || opts.smoothness < 0 {
		return nil, errors.New("invalid smoothness value")
	}

	passes := 8
	scale := 0.01 * (1 - opts.smoothness)
	nt := data.PerlinNoise(data.PerlinNoiseOptions{
		X:      opts.x,
		Y:      opts.y,
		Seed:   opts.seed,
		Size:   opts.size,
		Scale:  scale,
		Passes: passes,
	})

	return nt, nil
}

func WriteBmp(img *image.RGBA, c *data.Chunk) {
	size := len(*c.Terrain)

	xStart := c.X
	yStart := c.Y

	for x := range size {
		for y := range size {
			point := (*c.Terrain)[x][y]
			c := point.Material.Color
			a := point.Elevation
			e := color.RGBA{
				R: a,
				G: a,
				B: a,
				A: 255,
			}
			c = color.RGBA(helpers.MergeColors(c, e, 0.5))

			img.Set(x+xStart, y+yStart, c)
		}
	}
}
