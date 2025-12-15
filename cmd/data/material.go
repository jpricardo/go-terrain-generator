package data

import "image/color"

type Material struct {
	name  string
	color color.RGBA
}

var (
	water = Material{name: "water", color: color.RGBA{R: 0, G: 48, B: 128, A: 255}}
	sand  = Material{name: "sand", color: color.RGBA{R: 212, G: 185, B: 104, A: 255}}
	grass = Material{name: "grass", color: color.RGBA{R: 48, G: 135, B: 15, A: 255}}
	stone = Material{name: "stone", color: color.RGBA{R: 69, G: 69, B: 69, A: 255}}
	snow  = Material{name: "snow", color: color.RGBA{R: 200, G: 200, B: 200, A: 255}}
)

func Water() *Material {
	return &water
}

func Sand() *Material {
	return &sand
}

func Grass() *Material {
	return &grass
}

func Stone() *Material {
	return &stone
}

func Snow() *Material {
	return &snow
}
