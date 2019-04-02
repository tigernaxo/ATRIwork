package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func main() {
	zr := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: 5000000,
			Y: 1,
		},
	}
	img := image.NewNRGBA(zr)
	resizedImg := resize.Resize(1920, 10, img, resize.Lanczos3)

	SavePng(resizedImg, "image.png")
}

// SavePng save image as png
func SavePng(img image.Image, name string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}
