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

// RGB is a subset of image/color.RGBA
// where we ignore the alpha channel
type RGB struct {
	R uint8
	G uint8
	B uint8
}

// NewRGB converts any color.Color to RGB.  Conversion may be lossy.
func NewRGB(c color.Color) RGB {
	c = color.RGBAModel.Convert(c)
	rgba, ok := c.(color.RGBA)
	if !ok {
		// This type assertion should always succeed
		// based on https://godoc.org/image/color#Model
		// If we see this panic triggered in practice we might need to revisit.
		panic("color conversion failure")
	}
	return RGB{
		R: rgba.R,
		G: rgba.G,
		B: rgba.B,
	}
}

// String formats color as six digit hexidecimal
func (c RGB) String() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}
