package data

import (
	"fmt"
)

type Filter struct {
	name    string
	ammount float64
	radius  int
}

func Contrast(ammount float64) Filter {

	return Filter{name: "contrast", ammount: ammount}
}

func applyContrastTo(t *Texture, ammount float64) *Texture {
	size := len(*t)
	var mean uint8 // TODO - median?
	var max uint8 = 0
	var min uint8 = 255
	total := 0

	for x := range size {
		for y := range size {
			v := (*t)[x][y]
			if v > max {
				max = v
			}

			if v < min {
				min = v
			}

			total += int((*t)[x][y])
		}
	}

	mean = uint8(total / (size * size))

	for x := range size {
		for y := range size {
			v := (*t)[x][y]
			nv := v

			if v > mean {
				d := 1 - (float64(v) / float64(max))
				c := d * ammount * (255 - float64(max))
				nv = v + uint8(c)
			} else {
				d := (float64(v) / float64(min))
				c := d * ammount * -float64(min)
				nv = v + uint8(c)
			}

			(*t)[x][y] = nv
		}
	}

	return t
}

func Blur(ammount float64, radius int) Filter {

	return Filter{name: "blur", ammount: ammount, radius: radius}
}

func applyBlurTo(t *Texture, ammount float64, radius int) *Texture {
	size := len(*t)

	for x := range size {
		for y := range size {
			total := 0
			count := 0

			for rx := x - radius; rx < x+radius; rx++ {
				for ry := y - radius; ry < y+radius; ry++ {
					if rx < 0 || rx >= size {
						continue
					}

					if ry < 0 || ry >= size {
						continue
					}

					total += int((*t)[rx][ry])
					count++
				}
			}

			mean := total / count
			c := (*t)[x][y]
			d := 1 - float64(c)/float64(mean)
			vc := d * float64(mean) * ammount
			nv := uint8(float64(c) + vc)

			(*t)[x][y] = nv
		}
	}

	return t
}

func (f *Filter) ApplyTo(t *Texture) *Texture {
	fmt.Printf("[FILTER] Applying filter %+v ...\n", *f)

	switch f.name {
	case "contrast":
		applyContrastTo(t, f.ammount)

	case "blur":
		applyBlurTo(t, f.ammount, f.radius)
	}

	return t
}
