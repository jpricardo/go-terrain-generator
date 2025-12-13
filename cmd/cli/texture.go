package main

import (
	"errors"
	"fmt"
	"image"
	"math"
	"math/rand/v2"
	"slices"
)

type TextureOptions struct {
	size       int
	baseLine   uint8
	smoothness float64
}

type Texture [][]uint8

func (t *Texture) whiteNoise(opts TextureOptions) (Texture, error) {
	if opts.smoothness > 1 || opts.smoothness < 0 {
		return nil, errors.New("invalid smoothness value")
	}

	fmt.Printf("[TEXTURE] Generating %dx%d white noise texture with %+v...\n", opts.size, opts.size, opts)

	texture := Texture{}

	for x := range opts.size {
		texture = append(texture, make([]uint8, opts.size))

		for y := range opts.size {
			r := uint8(rand.Uint64()%255) + 1
			wd := 1 - float64(r)/float64(opts.baseLine)
			wc := wd * float64(opts.baseLine) * math.Pow(opts.smoothness, 2)
			e := float64(r) + wc

			texture[x][y] = uint8(math.Min(e, float64(maxElevation)))
		}
	}

	return texture, nil
}

type MergeOptions struct {
	smoothness float64
	opacity    float64
}

func (t *Texture) merge(textures []Texture, opts MergeOptions) (Texture, error) {
	slices.SortFunc(textures, func(i, j Texture) int {
		return len(i) - len(j)
	})

	if len(textures[0]) > len(*t) {
		return nil, errors.New("merged textures should be smaller than the target texture")
	}

	for _, texture := range textures {
		ts := len(texture)
		p := len(*t) / len(texture)

		for x := range len(*t) {
			for y := range len(*t) {
				c := (*t)[x][y]

				tx := int(math.Min(float64(x/p), float64(ts-1)))
				ty := int(math.Min(float64(y/p), float64(ts-1)))

				v := texture[tx][ty]
				d := 1 - float64(c)/float64(v)
				vc := d * float64(v) * opts.smoothness * opts.opacity
				nv := uint8(float64(c) + vc)

				(*t)[x][y] = nv
			}
		}
	}

	return *t, nil
}

type ApplyOptions struct {
	smoothness float64
	opacity    float64
}

func (t *Texture) toBitmap() *image.Gray {
	size := len(*t)
	img := image.NewGray(image.Rect(0, 0, size, size))

	for y := range size {
		for x := range size {
			offset := (y * img.Stride) + x
			img.Pix[offset] = (*t)[y][x]
		}
	}

	return img
}
