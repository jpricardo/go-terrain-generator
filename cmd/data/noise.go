package data

import (
	"errors"
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
	X      int
	Y      int
	Seed   uint64
	Size   int
	Scale  float64
	Passes int
}

func PerlinNoise(opts PerlinNoiseOptions) *Texture {
	p := newPerlin(opts.Seed)
	t := Texture{}

	for x := range opts.Size {
		tx := x + opts.X
		t = append(t, make([]uint8, opts.Size))

		for y := range opts.Size {
			ty := y + opts.Y

			var noiseSum float64
			amplitude := 16.0 * opts.Scale
			frequency := 0.0002 * (1 - opts.Scale)
			maxValue := 0.0

			for i := 0; i < opts.Passes; i++ {
				noiseSum += p.getNoise(float64(tx)*frequency, float64(ty)*frequency) * amplitude
				maxValue += amplitude

				if i == 0 {
					amplitude = 2
					frequency = opts.Scale
				} else {
					amplitude *= 0.25
					frequency *= 4
				}
			}

			normalized := (noiseSum/maxValue + 1.0) / 2.0

			t[x][y] = uint8(normalized * 255)
		}
	}

	return &t
}

func newPerlin(seed uint64) *Perlin {
	p := &Perlin{}
	r := rand.New(rand.NewPCG(seed, 0))

	var perm [256]int
	for i := 0; i < 256; i++ {
		perm[i] = i
	}

	for i := 0; i < 256; i++ {
		swap := uint8(r.Uint64()%255) + 1
		perm[i], perm[swap] = perm[swap], perm[i]
	}

	for i := 0; i < 256; i++ {
		p.p[i] = perm[i]
		p.p[i+256] = perm[i]
	}
	return p
}

func (p Perlin) getNoise(x, y float64) float64 {
	X := int(math.Floor(x)) & 255
	Y := int(math.Floor(y)) & 255

	x -= math.Floor(x)
	y -= math.Floor(y)

	u := helpers.Fade(x)
	v := helpers.Fade(y)

	a := p.p[X] + Y
	aa := p.p[a]
	ab := p.p[a+1]
	b := p.p[X+1] + Y
	ba := p.p[b]
	bb := p.p[b+1]

	return helpers.Lerp(v,
		helpers.Lerp(u, helpers.Grad(p.p[aa], x, y), helpers.Grad(p.p[ba], x-1, y)),
		helpers.Lerp(u, helpers.Grad(p.p[ab], x, y-1), helpers.Grad(p.p[bb], x-1, y-1)))
}
