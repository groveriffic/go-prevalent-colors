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

	// color.RGBA keeps each channel as a uint8
	// if we ignore the A channel the is equivalent to 6 digit hexidecimal format
	colorModel := color.RGBAModel

	eachPixel(img, func(x, y int, c color.Color) {
		cc, ok := colorModel.Convert(c).(color.RGBA)
		if !ok {
			log.Fatal("color cast failed")
		}
		fmt.Println(x, y, cc.R, cc.G, cc.B, cc.A)
	})
}

func eachPixel(img image.Image, f func(int, int, color.Color)) {
	r := img.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			f(x, y, img.At(x, y))
		}
	}
}
