package color

import (
	"fmt"
	"image"
	"image/color"
	"sort"
)

/* OPTIMIZE: Suspect image.Image uses a fair bit of space. Ideally I would
     like an API that just streams out pixels without needing to keep the
		 full image in memory.
		 Factors to consider:
		  - Is this available or even possible with all image encodings?
			- Need benchmarking to verify improvement
*/
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

// Counter is used to tally all colors in an image
type Counter map[RGB]int

/* OPTIMIZE: I would be tempted to try a probablistic data structure here to
    reduce space usage.  Count-min sketch might be interesting.
		Factors to consider:
		 - Need benchmarking to verify this is the biggest bottleneck
		 - Need benchmarking to verify improvement
		 - Need to be able to list keys (not sure of this with Count-min sketch)
		 - Accuracy would be lost (Suspect that is acceptible here)
*/

// Inc increments a single color's count
func (cc Counter) Inc(c RGB) {
	i := cc[c]
	cc[c] = i + 1
}

// Image counts the occurances of each RGB color in an image
func (cc Counter) Image(img image.Image) {
	eachPixel(img, func(x int, y int, c color.Color) {
		rgb := NewRGB(c)
		cc.Inc(rgb)
	})
}

// Rank lists counted colors from most occurances to least
func (cc Counter) Rank() []RGB {
	colors := []RGB{}
	for rgb := range cc {
		colors = append(colors, rgb)
	}
	less := func(i, j int) bool {
		iColor := colors[i]
		jColor := colors[j]
		iCount := cc[iColor]
		jCount := cc[jColor]
		return iCount < jCount
	}
	sort.Slice(colors, less)
	return colors
}

// TopThree returns the three most common colors
func (cc Counter) TopThree() []RGB {
	// OPTIMIZE: This could be faster by not sorting all colors and only keeping
	// The top 3.
	// https://en.wikipedia.org/wiki/Partial_sorting
	colors := cc.Rank()
	return colors[:3]
}
