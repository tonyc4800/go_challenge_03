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

// calcAvgRGBm accepts and image and returns the average pixel values for each
// channel as an 8-bit float64 array.
func calcAvgRGB(img image.Image) [3]float64 {
	bounds := img.Bounds()

	var rgbVals [3]float64
	var totalPix float64

	// Loop image from bottom left to upper right.  Values are shifted by 8
	// since RGBA returns values on [0, 65535](16-bit) and [0, 255](8-bit) is,
	// subjectively, easier to interpret.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgbVals[0] = rgbVals[0] + float64(r>>8)
			rgbVals[1] = rgbVals[1] + float64(g>>8)
			rgbVals[2] = rgbVals[2] + float64(b>>8)
			totalPix++
		}
	}

	// Calculate average for each channel.
	rgbVals[0] = rgbVals[0] / totalPix
	rgbVals[1] = rgbVals[1] / totalPix
	rgbVals[2] = rgbVals[2] / totalPix

	return rgbVals
}

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

	AvgRGB := calcAvgRGB(img)

	fmt.Printf("%-8s %-8s %-8s\n", "red", "green", "blue")
	fmt.Printf("%6.2f %6.2f %6.2f\n", AvgRGB[0], AvgRGB[1], AvgRGB[2])

	fmt.Println("yipee")
}
