package helpers

import "image/color"

func MergeColors(c1, c2 color.RGBA, opacity float64) color.RGBA {

	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	r_merged := (r1 + uint32(float64(r2)*opacity)) / 2
	g_merged := (g1 + uint32(float64(g2)*opacity)) / 2
	b_merged := (b1 + uint32(float64(b2)*opacity)) / 2
	a_merged := (a1 + a2) / 2

	return color.RGBA{
		R: uint8(r_merged >> 8),
		G: uint8(g_merged >> 8),
		B: uint8(b_merged >> 8),
		A: uint8(a_merged >> 8),
	}

}
