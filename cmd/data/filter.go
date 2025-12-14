package data

import (
	"fmt"
)

type Filter struct {
	name    string
	ammount float64
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
				d := float64(1 - (v / max))
				c := d * ammount * (255 - float64(max))
				nv = v + uint8(c)
			}
			if v < mean {
				d := float64((v / min))
				c := d * ammount * -float64(min)
				nv = v + uint8(c)
			}

			(*t)[x][y] = nv
		}
	}

	return t
}

func (f *Filter) ApplyTo(t *Texture) *Texture {
	fmt.Printf("[FILTER] Applying %f %s filter...\n", f.ammount, f.name)

	switch f.name {
	case "contrast":
		applyContrastTo(t, f.ammount)
	}

	return t
}
