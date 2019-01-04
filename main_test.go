package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"os"
	"testing"
)

func TestEachPixel(t *testing.T) {
	var i int

	f, err := os.Open("fixtures/10x10.jpg")
	if err != nil {
		t.Error(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		t.Error(err)
	}

	err = f.Close()
	if err != nil {
		t.Error(err)
	}

	eachPixel(img, func(x, y int, c color.Color) {
		i++
	})

	t.Log(i)
	if i != 100 {
		t.Error("expected function to be called exactly 100 times")
	}
}

func ExampleRGB_String() {
	c := RGB{R: 0, G: 0, B: 0}
	fmt.Println(c.String())

	c = RGB{R: 0, G: 255, B: 0}
	fmt.Println(c.String())

	c = RGB{R: 255, G: 255, B: 255}
	fmt.Println(c.String())

	// Output:
	// #000000
	// #00FF00
	// #FFFFFF
}

var _ fmt.Stringer = RGB{} // Interface check

func ExampleNewRGB() {
	rgb := NewRGB(color.Gray{Y: 0})
	fmt.Println(rgb.String())

	rgb = NewRGB(color.Gray{Y: 255})
	fmt.Println(rgb.String())

	// YCbCr is color model used by JPEG
	rgb = NewRGB(color.YCbCr{Y: 255, Cb: 255, Cr: 255})
	fmt.Println(rgb.String())

	// Output:
	// #000000
	// #FFFFFF
	// #FF79FF
}
