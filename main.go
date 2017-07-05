package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
)

// image documentation: https://golang.org/pkg/image/

// ----------- Mosaic Image
// Read image

// Determine avg color

// Break into color "buckets"

// Break each bucket into Light - Med - Dark (or range, based on defined value)

// ---------- Target Image
// Read target image

// Downsample target image to X,Y

// Determine average pix color

// Break colors into "buckets" [determined by `mosaic buckets`]

// Light - Med - Dark color for each bucket (or range, based on defined value)

// ------------ Image Creation
// Map Mosaic images to target image

// Create new image w/each mosaic image mapped to the target image as a 'pixel'

func main() {
	// Read mosaic images to see what values we have to work with.

	// Read in target image to see how we have to map

	// Map mosaic images to target image values

	// Create final image

	// Profit

	tarImgP := "./input/target/day_man.jpg"
	reader, err := os.Open(tarImgP)
	if err != nil {
		fmt.Println("can't open img")
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		fmt.Println("can't read img")
	}
	bounds := img.Bounds()

	var hist [16][3]int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// values in range [0, 65535].
			// Shifting by 12, reduces range to [0,15].
			hist[r>>12][0]++
			hist[g>>12][1]++
			hist[b>>12][2]++
		}
	}

	fmt.Printf("%-14s %6s %6s %6s\n", "bin", "red", "green", "blue")
	for i, x := range hist {
		fmt.Printf("0x%04x-0x%04x: %6d %6d %6d\n", i<<12, (i+1)<<12-1, x[0], x[1], x[2])
	}

	fmt.Println("yipee")
}
