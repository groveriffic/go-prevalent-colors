package color

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"os"
	"testing"
)

var fixtureImageA image.Image

func TestMain(m *testing.M) {
	f, err := os.Open("../fixtures/10x10.jpg")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	fixtureImageA = img

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestEachPixel(t *testing.T) {
	var i int

	eachPixel(fixtureImageA, func(x, y int, c color.Color) {
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

func TestColorCounter_Inc(t *testing.T) {
	cc := make(ColorCounter)
	c := RGB{}
	cc.Inc(c)
	i, ok := cc[c]
	if !ok {
		t.Error("incremented color not in map")
	}
	if i != 1 {
		t.Error("incremented color count incorrect")
	}
}

func TestColorCounter_Image(t *testing.T) {
	cc := make(ColorCounter)
	cc.Image(fixtureImageA)

	rgb := RGB{R: 0x5c, G: 0x6c, B: 0x83}
	i, ok := cc[rgb]
	if !ok {
		t.Error("expected color not in map")
	}
	if i != 2 {
		t.Error("expected color count incorrect")
	}
}

func TestColorCounter_Rank(t *testing.T) {
	cc := ColorCounter{}
	for i := 0; i < 5; i++ {
		rgb := RGB{R: uint8(i), G: uint8(i), B: uint8(i)}
		cc[rgb] = i
	}

	colors := cc.Rank()
	t.Logf("%#v\n", colors)
	if colors[0].R != 0 {
		t.Fail()
	}
	if colors[1].R != 1 {
		t.Fail()
	}
	if colors[2].R != 2 {
		t.Fail()
	}
	if colors[3].R != 3 {
		t.Fail()
	}
	if colors[4].R != 4 {
		t.Fail()
	}
}

func TestColorCounter_TopThree(t *testing.T) {
	cc := ColorCounter{}
	for i := 0; i < 5; i++ {
		rgb := RGB{R: uint8(i), G: uint8(i), B: uint8(i)}
		cc[rgb] = i
	}

	colors := cc.TopThree()
	t.Logf("%#v\n", colors)
	if colors[0].R != 0 {
		t.Fail()
	}
	if colors[1].R != 1 {
		t.Fail()
	}
	if colors[2].R != 2 {
		t.Fail()
	}
	if len(colors) != 3 {
		t.Fail()
	}
}
