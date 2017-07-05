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

	// float64 used since we need a floating point.
	var rgbVal [3]float64
	var totalPixCount float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgbVal[0] = rgbVal[0] + float64(r)
			rgbVal[1] = rgbVal[1] + float64(g)
			rgbVal[2] = rgbVal[2] + float64(b)
			totalPixCount++
		}
	}

	// calculate average
	rgbVal[0] = rgbVal[0] / totalPixCount
	rgbVal[1] = rgbVal[1] / totalPixCount
	rgbVal[2] = rgbVal[2] / totalPixCount

	fmt.Printf("%-8s %-8s %-8s\n", "red", "green", "blue")
	fmt.Printf("%6.2f %6.2f %6.2f\n", rgbVal[0], rgbVal[1], rgbVal[2])

	fmt.Println("yipee")
}
