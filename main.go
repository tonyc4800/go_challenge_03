package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
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
func calcAvgRGB(img image.Image) [3]uint32 {
	bounds := img.Bounds()

	rgbVals := [3]uint32{0, 0, 0}
	var totalPix uint32

	// Loop image from bottom left to upper right.  Values are shifted by 8
	// since RGBA returns values on [0, 65535](16-bit) and [0, 255](8-bit) is,
	// subjectively, easier to interpret.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			//fmt.Printf("%v:{%v} %v:{%v} %v:{%v}\n", "red", r/256, "green", g/256, "blue", b/256)
			rgbVals[0] = rgbVals[0] + (r / 256)
			rgbVals[1] = rgbVals[1] + (g / 256)
			rgbVals[2] = rgbVals[2] + (b / 256)
			totalPix++
		}
	}
	//fmt.Printf("%v:{%v} %v:{%v} %v:{%v}\n", "red", rgbVals[0], "green", rgbVals[1], "blue", rgbVals[2])

	// Calculate average for each channel.
	rgbVals[0] = rgbVals[0] / totalPix
	rgbVals[1] = rgbVals[1] / totalPix
	rgbVals[2] = rgbVals[2] / totalPix

	return rgbVals
}

// TODO: complete this function.
// resizeImage accepts and image and target x and y sizes, then resizes and
// returns the image. Docs: https://golang.org/pkg/image/#NewRGBA
func resizeImage(oImg image.Image, tWidth int, tHeight int) image.Image {
	// TODO: Ensure target size is under original size.

	// Create new, resized, image rectangle.
	rImage := image.NewRGBA(image.Rect(0, 0, tWidth, tHeight))

	bounds := oImg.Bounds()
	oWidth := bounds.Max.X - bounds.Min.X
	oHeight := bounds.Max.Y - bounds.Min.Y
	wRatio := float64(oWidth) / float64(tWidth)
	hRatio := float64(oHeight) / float64(tHeight)

	// Create a grid of coordinates for subimages from the original image that
	// can be mapped into the new image.
	var xCoords []int
	var yCoords []int
	// fmt.Printf("wRatio: %v\n", wRatio)
	// fmt.Printf("hRatio: %v\n", hRatio)
	for y := 0; y < tHeight; y++ {

		// The coordinate value will be cropped to an int value, not rounded.
		i := int(float64(y) * hRatio)
		yCoords = append(yCoords, i)
	}
	for x := 0; x < tWidth; x++ {

		// The coordinate value will be cropped to an int value, not rounded.
		i := int(float64(x) * wRatio)
		xCoords = append(xCoords, i)
	}

	// Replace last value with max original value. NOTE: This will affect the
	// image quality on the right and upper edges.
	xCoords[len(xCoords)-1] = bounds.Max.X
	yCoords[len(yCoords)-1] = bounds.Max.Y

	// Remove first value from slice.
	xCoords = append(xCoords[:0], xCoords[0+1:]...)
	yCoords = append(yCoords[:0], yCoords[0+1:]...)

	fmt.Printf("Height: %v x %v\n", bounds.Min.Y, bounds.Max.Y)
	fmt.Printf("Width: %v x %v\n", bounds.Min.X, bounds.Max.X)
	// fmt.Println(xCoords)
	// fmt.Println(yCoords)

	// Loop coordinates and create sub images
	xStart := 0
	yStart := 0

	for j, yCoord := range yCoords {
		for i, xCoord := range xCoords {

			// (i, j) will be the coord for the pix value in the new image.
			// (xStart, yStart, xCoord, yCoord)Rect will be the sub image.
			subImg := image.NewRGBA(image.Rect(0, 0, xCoord-xStart, yCoord-yStart))

			// Fill subimage pixel values.
			n := 0

			fmt.Println("Subimage values")
			for yy := yStart; yy < yCoord; yy++ {

				m := 0
				for xx := xStart; xx < xCoord; xx++ {

					r, g, b, _ := oImg.At(xx, yy).RGBA()

					cVal := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
					fmt.Printf("{%v,%v}->{%v,%v}--%v:{%v} %v:{%v} %v:{%v}\n", xx, yy, m, n, "red", cVal.R, "green", cVal.G, "blue", cVal.B)
					subImg.SetRGBA(m, n, cVal)
					m++
				}
				n++
			}

			// Get average value
			imgVals := calcAvgRGB(subImg)

			// Assign value to new image. alpha is hardcoded to 255 since we do
			// not want a transparent image.
			nVal := color.RGBA{R: uint8(imgVals[0]), G: uint8(imgVals[1]), B: uint8(imgVals[2]), A: 255}
			fmt.Printf("(%v,%v)%v:{%v} %v:{%v} %v:{%v}\n", i, j, "red", nVal.R, "green", nVal.G, "blue", nVal.B)
			rImage.SetRGBA(i, j, nVal)

			// Update coordinate grid.
			xStart = xCoord
			fmt.Println("---------------")
		}
		yStart = yCoord
	}

	return rImage
}

func main() {
	// Read mosaic images to see what values we have to work with.
	// - read in
	// - downsample (resize)
	// - calculate avg for image
	// - create map mImgIndex:avgPixValue

	// Read in target image to see how we have to map
	// - read in
	// - downsample to target size (resize)
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

	// TODO: complete this function.
	resizedTargetImg := resizeImage(img, 200, 200)
	bounds := resizedTargetImg.Bounds()
	oWidth := bounds.Max.X - bounds.Min.X
	oHeight := bounds.Max.Y - bounds.Min.Y
	fmt.Printf("resizedTargetImg: %vx%v\n", oWidth, oHeight)

	rsImgF, err := os.Create("./output/resizedTarget.jpg")
	if err != nil {
		fmt.Printf("Error creating img file: %v\n", err)
	}

	err = jpeg.Encode(rsImgF, resizedTargetImg, nil)

	createdImgF := "./output/resizedTarget.jpg"
	readCreatedImg, err := returnImgFromPath(createdImgF)
	if err != nil {
		fmt.Printf("Error Obtaining Img: %v\n", err)
	}
	readImgBounds := readCreatedImg.Bounds()
	rsWidth := readImgBounds.Max.X - readImgBounds.Min.X
	rsHeight := readImgBounds.Max.Y - readImgBounds.Min.Y
	fmt.Printf("readCreatedImg: %vx%v\n", rsWidth, rsHeight)

	xx := 198
	yy := 198
	r, g, b, _ := resizedTargetImg.At(xx, yy).RGBA()
	fmt.Printf("{%v,%v}--%v:{%v} %v:{%v} %v:{%v}\n", xx, yy, "red", uint8(r), "green", uint8(g), "blue", uint8(b))
	r, g, b, _ = readCreatedImg.At(xx, yy).RGBA()
	fmt.Printf("{%v,%v}--%v:{%v} %v:{%v} %v:{%v}\n", xx, yy, "red", uint8(r), "green", uint8(g), "blue", uint8(b))

	fmt.Println("yipee")
}
