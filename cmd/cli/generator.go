package main

import (
	"errors"
	"fmt"
	"go-terrain-generator/cmd/data"
	"go-terrain-generator/cmd/helpers"
	"log"
)

type ChunkOptions struct {
	x          int
	y          int
	seed       uint64
	size       int
	smoothness float64
}

func GenerateChunk(ch chan *data.Chunk, opts ChunkOptions) {
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

	img := ct.ToBitmap()
	helpers.SaveBmp(img, outputDir, fmt.Sprintf("chunk_%d_%d.png", opts.x, opts.y))

	ch <- c
}

type TerrainOptions struct {
	seed       uint64
	size       int
	smoothness float64
}

func GenerateTerrain(opts TerrainOptions) (*data.Terrain, error) {
	if opts.size <= 0 || opts.size > data.MaxSize || opts.size < data.MinSize {
		return nil, errors.New("invalid map size")
	}

	t := data.NewTerrain(opts.size)

	chunkSize := data.MinSize
	chunks := opts.size / chunkSize
	ch := make(chan *data.Chunk)

	for x := range chunks {
		for y := range chunks {
			go GenerateChunk(ch, ChunkOptions{
				x:          x,
				y:          y,
				seed:       opts.seed,
				smoothness: opts.smoothness,
				size:       chunkSize,
			})
		}
	}

	for range chunks * chunks {
		c := <-ch
		t = t.ApplyChunk(c)
	}

	img := t.ToBitmap()
	helpers.SaveBmp(img, outputDir, "chunked_terrain.png")

	return t, nil
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
