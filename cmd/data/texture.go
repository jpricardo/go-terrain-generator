package data

import (
	"errors"
	"image"
	"slices"
)

type Texture [][]uint8

type MergeOptions struct {
	opacity float64
}

func (t *Texture) Merge(textures []Texture, opts MergeOptions) (*Texture, error) {
	slices.SortFunc(textures, func(i, j Texture) int {
		return len(i) - len(j)
	})

	max := int(0)

	for x := range len(*t) {
		for y := range len(*t) {
			total := int(0)

			for _, texture := range textures {
				if len(texture) != len(*t) {
					return nil, errors.New("textures should be of the same size")
				}
				total += int(texture[x][y])
			}

			if total > max {
				max = total
			}

			c := (*t)[x][y]
			mean := total / len(textures)
			d := 1 - float64(c)/float64(mean)
			vc := d * float64(mean) * opts.opacity
			nv := uint8(float64(c) + vc)

			(*t)[x][y] = nv
		}
	}

	return t, nil
}

func (t *Texture) ApplyFilters(filters []Filter) *Texture {

	for _, f := range filters {
		t = f.ApplyTo(t)
	}

	return t
}

func (t *Texture) ToBitmap() *image.Gray {
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
