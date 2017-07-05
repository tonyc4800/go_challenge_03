package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
)

// image documentation: https://golang.org/pkg/image/

// returnImgFromPath accepts a file path to a jpeg image and returns the image.
func returnImgFromPath(imgPath string) (image.Image, error) {
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open img: %v", err)
	}
	defer f.Close()

	//img, _, err := image.Decode(r)
	//img, err := jpeg.Decode(f)
	img, err := png.Decode(f)
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

	// Add max original value to the slice.
	xCoords = append(xCoords, bounds.Max.X)
	yCoords = append(yCoords, bounds.Max.Y)

	// Remove first value from slice.
	xCoords = append(xCoords[:0], xCoords[0+1:]...)
	yCoords = append(yCoords[:0], yCoords[0+1:]...)

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
			for yy := yStart; yy <= yCoord; yy++ {
				m := 0
				for xx := xStart; xx <= xCoord; xx++ {
					r, g, b, _ := oImg.At(xx, yy).RGBA()
					cVal := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
					subImg.Set(m, n, cVal)
					m++
				}
				n++
			}

			// Assign value to new image. alpha is hardcoded to 255 since we do
			// not want a transparent image.
			imgVals := calcAvgRGB(subImg)
			nVal := color.RGBA{R: uint8(imgVals[0]), G: uint8(imgVals[1]), B: uint8(imgVals[2]), A: 255}
			rImage.Set(i, j, nVal)

			xStart = xCoord
		}
		yStart = yCoord
	}

	return rImage
}

// TODO: create this function
// nearestMapping...
func nearestMapping() {

}

// writeImgToFile is a convenience function that will likely be deleted.
func writeImgToFile(img image.Image, filePath string) error {
	//rsImgF, err := os.Create()
	rsImgF, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating img file: %v\n", err)
	}
	//defer rsImgF.Close()
	err = png.Encode(rsImgF, img)
	return err
}

func main() {

	// Read mosaic images to see what values we have to work with.
	mosDir := "./input/mosaic/PCB_square_png"

	// create directory to hold smaller images (if not exist) 777
	smallPath := mosDir + "/resized"
	if _, err := os.Stat(smallPath); os.IsNotExist(err) {
		os.Mkdir(smallPath, os.ModePerm)
	}

	oFiles, _ := ioutil.ReadDir(mosDir)
	//sFiles, _ := ioutil.ReadDir(smallPath)
	//if len(oFiles) != len(sFiles) {
	for _, f := range oFiles {
		fPath := f.Name()
		ext := filepath.Ext(fPath)
		key := fPath[0 : len(fPath)-len(ext)]

		if ext == ".png" || ext == ".jpg" {
			img, err := returnImgFromPath(mosDir + "/" + fPath)
			if err != nil {
				fmt.Println(fPath)
				fmt.Printf("Error Obtaining Img: %v\n", err)
			}

			rsImg := resizeImage(img, 60, 60)
			rsPath := smallPath + "/" + fPath
			writeImgToFile(rsImg, rsPath)

			imgVals := calcAvgRGB(rsImg)

			//fmt.Println(key)
			fmt.Printf("%-15s (r:%v,g:%v,b:%v)\n", key, uint8(imgVals[0]), uint8(imgVals[1]), uint8(imgVals[2]))
		}

	}
	//}

	// X read in
	// X downsample (resize)
	// X calculate avg for image
	// - create map mImgIndex:avgPixValue

	// Read in target image to see how we have to map
	// X read in
	// X downsample to target size (resize)
	// - create map pixIndex:avgPixValue

	// Image Creation
	// - map mosaic images to target image pixels
	// - create new image

	// TODO: w/in resize, make sure we write to the final col/row

	// TODO: Create nearest mapping function to map pixel value to nearest
	// mosaic value.

	// Profit

	// tarImgP := "./input/target/day_man.png"

	// img, err := returnImgFromPath(tarImgP)
	// if err != nil {
	// 	fmt.Printf("Error Obtaining Img: %v\n", err)
	// }

	// resizedTargetImg := resizeImage(img, 200, 200)

	// err = writeImgToFile(resizedTargetImg, "./output/resizedTarget.png")

	fmt.Println("yipee")
}
