package helpers

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func SaveBmp(img *image.Gray, outputDir string, fileName string) {
	f, err := os.Create(fmt.Sprintf("%s%s", outputDir, fileName))
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Panic(err)
	}
}
