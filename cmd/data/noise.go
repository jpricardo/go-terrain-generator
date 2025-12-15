package data

import (
	"errors"
	"fmt"
	"go-terrain-generator/cmd/helpers"
	"math"
	"math/rand/v2"
)

type WhiteNoiseOptions struct {
	size       int
	baseLine   uint8
	smoothness float64
}

func WhiteNoise(opts WhiteNoiseOptions) (*Texture, error) {
	if opts.smoothness > 1 || opts.smoothness < 0 {
		return nil, errors.New("invalid smoothness value")
	}

	fmt.Printf("[TEXTURE] Generating white noise texture with %+v...\n", opts)

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

	return &texture, nil
}

type ValueNoiseOptions struct {
	size       int
	baseLine   uint8
	smoothness float64
}

type Perlin struct {
	p [512]int
}

type PerlinNoiseOptions struct {
	Size   int
	Scale  float64
	Passes int
}

func PerlinNoise(opts PerlinNoiseOptions) *Texture {
	p := newPerlin()
	t := Texture{}

	fmt.Printf("[TEXTURE] Generating Perlin noise texture with %+v...\n", opts)

	for x := range opts.Size {
		t = append(t, make([]uint8, opts.Size))

		for y := range opts.Size {

			var noiseSum float64
			amplitude := 1.0
			frequency := opts.Scale
			maxValue := 0.0

			for i := 0; i < opts.Passes; i++ {
				noiseSum += p.getNoise(float64(x)*frequency, float64(y)*frequency) * amplitude
				maxValue += amplitude

				amplitude *= 0.25
				frequency *= 2
			}

			normalized := (noiseSum/maxValue + 1.0) / 2.0

			t[x][y] = uint8(normalized * 255)
		}
	}

	return &t
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
