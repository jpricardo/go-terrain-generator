package main

import (
	"errors"
	"fmt"
	"go-terrain-generator/cmd/data"
	"go-terrain-generator/cmd/helpers"
)

type TerrainOptions struct {
	size       int
	smoothness float64
}

func GenerateTerrain(opts TerrainOptions) (*data.Terrain, error) {
	if opts.size <= 0 || opts.size > data.MaxSize || opts.size < data.MinSize {
		return nil, errors.New("invalid map size")
	}

	fmt.Printf("[TERRAIN] Generating terrain map with %+v...\n", opts)

	t := data.NewTerrain(opts.size)
	et, err := GenerateElevation(ElevationOptions{smoothness: opts.smoothness, size: opts.size})
	if err != nil {
		return nil, err
	}

	t, err = t.ApplyElevation(et, data.ApplyOptions{Opacity: 0.5 - (0.25 * opts.smoothness)})
	if err != nil {
		return nil, err
	}

	return t, nil
}

type ElevationOptions struct {
	size       int
	smoothness float64
}

func GenerateElevation(opts ElevationOptions) (*data.Texture, error) {
	if opts.smoothness > 1 || opts.smoothness < 0 {
		return nil, errors.New("invalid smoothness value")
	}

	passes := int(4 / opts.smoothness)
	scale := (1 - opts.smoothness) * data.MinSize / float64(opts.size)
	nt := data.PerlinNoise(data.PerlinNoiseOptions{
		Size:   opts.size,
		Scale:  scale,
		Passes: passes,
	})

	img := nt.ToBitmap()
	helpers.SaveBmp(img, outputDir, "elevation_texture.png")

	return nt, nil
}
