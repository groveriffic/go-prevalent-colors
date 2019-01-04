package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"os"
)

func main() {
	filename := "FApqk3D.jpg"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// m := make(map[color.Color]int)
	b := img.Bounds()

	// color.RGBA keeps each channel as a uint8
	// if we ignore the A channel the is equivalent to 6 digit hexidecimal format
	colorModel := color.RGBAModel

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c, ok := colorModel.Convert(img.At(x, y)).(color.RGBA)
			if !ok {
				log.Fatal("color cast failed")
			}
			fmt.Println(x, y, c.R, c.G, c.B, c.A)
		}
	}
}
