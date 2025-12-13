package main

import (
	"errors"
	"fmt"
	"image"
	"math"
	"math/rand/v2"
)

type Terrain [][]*Point

type Point struct {
	elevation uint8
}

const (
	waterLevel   = uint8(64)
	maxElevation = uint8(255)
	lowResRatio  = float64(0.25)
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

type ElevationOptions struct {
	smoothness float64
}

func (t *Terrain) applyLayer(layer Terrain, opts ElevationOptions) (Terrain, error) {
	ts := len(*t)
	ls := len(layer)

	if len(layer) <= 0 {
		return nil, errors.New("invalid layer size")
	}

	fmt.Printf("Applying %dx%d layer to %dx%d terrain...\n", ls, ls, ts, ts)

	p := ts / ls
	for x := range ts {
		row := (*t)[x]

		for y := range row {
			c := row[y].elevation

			lx := int(math.Min(float64(x/p), float64(ls-1)))
			ly := int(math.Min(float64(y/p), float64(ls-1)))

			te := layer[lx][ly].elevation
			td := 1 - float64(c)/float64(te)
			tc := td * float64(te) * opts.smoothness
			e := float64(c) + tc

			row[y].elevation = uint8(math.Min(e, float64(maxElevation)))
		}
	}

	return *t, nil
}

func (t *Terrain) getBaseLayer(opts ElevationOptions) (Terrain, error) {
	if opts.smoothness > 1 || opts.smoothness < 0 {
		return nil, errors.New("invalid smoothness value")
	}

	size := len(*t)
	fmt.Printf("Generating %dx%d base layer with %+v...\n", size, size, opts)

	for x := range *t {
		row := (*t)[x]

		for y := range row {
			r := uint8(rand.Uint64() % 256)
			wd := 1 - float64(r)/float64(waterLevel)
			wc := wd * float64(waterLevel) * math.Pow(opts.smoothness, 2)
			e := float64(r) + wc
			row[y].elevation = uint8((math.Min(e, float64(maxElevation))))
		}
	}

	return *t, nil
}

func (t *Terrain) generateElevation(opts ElevationOptions) (Terrain, error) {
	if opts.smoothness > 1 || opts.smoothness < 0 {
		return nil, errors.New("invalid smoothness value")
	}

	passes := int(math.Max(1/lowResRatio, 1))
	fmt.Println("Total passes: ", passes)

	baseLayer, _ := t.getBaseLayer(opts)

	for p := 1; p <= passes; p++ {
		ls := float64(p) / lowResRatio
		lt := initializeTerrain(int(ls))
		lt, _ = lt.getBaseLayer(opts)
		baseLayer.applyLayer(lt, ElevationOptions{smoothness: math.Pow(opts.smoothness, 8)})
	}

	return baseLayer, nil
}

func generateTerrain(size int) (Terrain, error) {
	if size <= 0 || size > maxSize {
		return nil, errors.New("invalid map size")
	}

	fmt.Printf("Generating %dx%d terrain map...\n", size, size)

	terrain := initializeTerrain(size)
	terrain, err := terrain.generateElevation(ElevationOptions{smoothness: 0.75})
	if err != nil {
		return nil, err
	}

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
