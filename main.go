package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
)

// image documentation: https://golang.org/pkg/image/

// returnImgFromPath accepts a file path to a jpeg image and returns the image.
func returnImgFromPath(imgPath string) (image.Image, error) {
	r, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open img: %v", err)
	}
	defer r.Close()

	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("unable read img: %v", err)
	}

	return img, nil
}

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

func main() {
	// Read mosaic images to see what values we have to work with.
	// - read in
	// - downsample
	// - calculate avg for image
	// - create map mImgIndex:avgPixValue

	// Read in target image to see how we have to map
	// - read in
	// - downsample to target size
	// - create map pixIndex:avgPixValue

	// Image Creation
	// - map mosaic images to target image pixels
	// - create new image

	// Profit

	tarImgP := "./input/target/day_man.jpg"

	img, err := returnImgFromPath(tarImgP)
	if err != nil {
		fmt.Printf("Error Obtaining Img: %v\n", err)
	}

	AvgRGB := calcAvgRGB(img)

	fmt.Printf("%-8s %-8s %-8s\n", "red", "green", "blue")
	fmt.Printf("%6.2f %6.2f %6.2f\n", AvgRGB[0], AvgRGB[1], AvgRGB[2])

	fmt.Println("yipee")
}
