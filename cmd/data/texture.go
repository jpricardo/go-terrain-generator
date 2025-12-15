package data

import (
	"errors"
	"fmt"
	"go-terrain-generator/cmd/helpers"
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

func (t *Texture) WhiteNoise(opts TextureOptions) (Texture, error) {
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

			texture[x][y] = uint8(e)
		}
	}

	return texture, nil
}

func (t *Texture) ValueNoise(opts TextureOptions) (Texture, error) {
	if opts.smoothness > 1 || opts.smoothness < 0 {
		return nil, errors.New("invalid smoothness value")
	}

	fmt.Printf("[TEXTURE] Generating %dx%d value noise texture with %+v...\n", opts.size, opts.size, opts)

	texture := Texture{}
	for x := range opts.size {
		texture = append(texture, make([]uint8, opts.size))

		for y := range opts.size {
			r := (helpers.HashInt(x) + helpers.HashInt(y)) / 2
			wd := 1 - float64(r)/float64(opts.baseLine)
			wc := wd * float64(opts.baseLine) * math.Pow(opts.smoothness, 2)
			e := float64(r) + wc

			texture[x][y] = uint8(e)
		}
	}

	return texture, nil
}

type Perlin struct {
	p [512]int
}

type PerlinNoiseOptions struct {
	size   int
	scale  float64
	passes int
}

func PerlinNoise(opts PerlinNoiseOptions) Texture {
	p := newPerlin()
	t := Texture{}

	fmt.Printf("[TEXTURE] Generating %dx%d Perlin noise texture with %+v...\n", opts.size, opts.size, opts)

	for x := range opts.size {
		t = append(t, make([]uint8, opts.size))

		for y := range opts.size {

			var noiseSum float64
			amplitude := 1.0
			frequency := opts.scale
			maxValue := 0.0

			for i := 0; i < opts.passes; i++ {
				noiseSum += p.getNoise(float64(x)*frequency, float64(y)*frequency) * amplitude
				maxValue += amplitude

				amplitude *= 0.25
				frequency *= 2
			}

			normalized := (noiseSum/maxValue + 1.0) / 2.0

			t[x][y] = uint8(normalized * 255)
		}
	}

	return t
}

func newPerlin() *Perlin {
	p := &Perlin{}
	// r := rand.New(rand.Source(seed))

	// Preenche a primeira metade com 0..255
	var perm [256]int
	for i := 0; i < 256; i++ {
		perm[i] = i
	}

	// Embaralha
	for i := 0; i < 256; i++ {
		swap := uint8(rand.Uint64()%255) + 1
		perm[i], perm[swap] = perm[swap], perm[i]
	}

	// Duplica para evitar overflow de índice
	for i := 0; i < 256; i++ {
		p.p[i] = perm[i]
		p.p[i+256] = perm[i]
	}
	return p
}

func (p Perlin) getNoise(x, y float64) float64 {
	// Encontra a célula unitária que contém o ponto
	X := int(math.Floor(x)) & 255
	Y := int(math.Floor(y)) & 255

	// Coordenadas relativas dentro da célula (0.0 a 1.0)
	x -= math.Floor(x)
	y -= math.Floor(y)

	// Calcula curvas de fade
	u := helpers.Fade(x)
	v := helpers.Fade(y)

	// Hash das coordenadas dos 4 cantos do quadrado
	a := p.p[X] + Y
	aa := p.p[a]
	ab := p.p[a+1]
	b := p.p[X+1] + Y
	ba := p.p[b]
	bb := p.p[b+1]

	// Interpolação linear entre os gradientes dos cantos
	return helpers.Lerp(v,
		helpers.Lerp(u, helpers.Grad(p.p[aa], x, y), helpers.Grad(p.p[ba], x-1, y)),
		helpers.Lerp(u, helpers.Grad(p.p[ab], x, y-1), helpers.Grad(p.p[bb], x-1, y-1)))
}

type MergeOptions struct {
	opacity float64
}

func (t *Texture) Merge(textures []Texture, opts MergeOptions) (Texture, error) {
	slices.SortFunc(textures, func(i, j Texture) int {
		return len(i) - len(j)
	})

	fmt.Printf("[TEXTURE] Merging %d textures with %+v\n", len(textures), opts)

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

	return *t, nil
}

func (t *Texture) ApplyFilters(filters []Filter) Texture {

	for _, f := range filters {
		t = f.ApplyTo(t)
	}

	return *t
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
