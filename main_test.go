package main

import (
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
