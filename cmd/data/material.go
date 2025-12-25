package data

import "image/color"

type Material struct {
	Name  string
	Color color.RGBA
}

var (
	water = Material{Name: "water", Color: color.RGBA{R: 0, G: 48, B: 128, A: 255}}
	sand  = Material{Name: "sand", Color: color.RGBA{R: 230, G: 200, B: 112, A: 255}}
	grass = Material{Name: "grass", Color: color.RGBA{R: 48, G: 135, B: 15, A: 255}}
	stone = Material{Name: "stone", Color: color.RGBA{R: 69, G: 69, B: 69, A: 255}}
	snow  = Material{Name: "snow", Color: color.RGBA{R: 200, G: 200, B: 200, A: 255}}
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
