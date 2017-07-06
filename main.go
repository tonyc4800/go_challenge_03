package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
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

func createMosaicMapping(mosDir string) map[string][3]uint8 {
	// create directory to hold smaller images (if not exist) 777
	smallPath := mosDir + "/resized"
	if _, err := os.Stat(smallPath); os.IsNotExist(err) {
		os.Mkdir(smallPath, os.ModePerm)
	}

	mosMap := make(map[string][3]uint8)
	oFiles, _ := ioutil.ReadDir(mosDir)

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

			rsImg := resizeImage(img, 35, 35)
			rsPath := smallPath + "/" + fPath
			writeImgToFile(rsImg, rsPath)

			imgVals := calcAvgRGB(rsImg)
			mVal := [3]uint8{uint8(imgVals[0]), uint8(imgVals[1]), uint8(imgVals[2])}
			mosMap[key] = mVal

		}
	}

	return mosMap
}

func main() {

	// Read mosaic images to see what values we have to work with.
	mosDir := "./input/mosaic/PCB_square_png"
	mosMap := createMosaicMapping(mosDir)

	// LOOK INTO: can a map be written to a file?

	tarImgP := "./input/target/day_man.png"

	img, err := returnImgFromPath(tarImgP)
	if err != nil {
		fmt.Printf("Error Obtaining Img: %v\n", err)
	}

	resizedTargetImg := resizeImage(img, 120, 120)

	bounds := resizedTargetImg.Bounds()
	rsWidth := bounds.Max.X - bounds.Min.X
	rsHeight := bounds.Max.Y - bounds.Min.Y

	// Loop resized image and map a mosaic value to the pixel value.
	mosKeyMap := [120 * 35][120 * 35]string{}
	track := 0
	for j := 0; j <= rsHeight; j++ {
		for i := 0; i <= rsWidth; i++ {
			r, g, b, _ := resizedTargetImg.At(i, j).RGBA()
			var mosaicN string
			closest := math.MaxFloat64
			for k, v := range mosMap {
				R := v[0]
				G := v[1]
				B := v[2]
				//fmt.Printf("%v:(R:%v, G:%v, B:%v)\n", k, R, G, B)

				// calculate nearest mosaic - weighted approach (since eyes are
				// more sensitive to G than B) -- 0.3, 0.59, 0.11 are magic
				// numbers referenced here:
				// `https://en.wikipedia.org/wiki/Luma_(video)`. The squareroot
				// is removed for optimization since we don't care what the
				// value of d is.
				rd := math.Pow((float64(R-uint8(r)) * 0.3), 2)
				gd := math.Pow((float64(G-uint8(g)) * 0.59), 2)
				bd := math.Pow((float64(B-uint8(b)) * 0.11), 2)
				d := rd + gd + bd
				if d < closest {
					closest = d
					mosaicN = k
				}

			}
			mosKeyMap[i][j] = mosaicN
			//fmt.Printf("%v, %v, (%v, %v, %v)->(%v, %v, %v)\n", track, mosaicN, mosMap[mosaicN][0], mosMap[mosaicN][1], mosMap[mosaicN][2], uint8(r), uint8(g), uint8(b))
			track++
		}
	}

	finalImage := image.NewRGBA(image.Rect(0, 0, 120*35, 120*35))
	fbounds := finalImage.Bounds()
	fWidth := fbounds.Max.X - fbounds.Min.X
	fHeight := fbounds.Max.Y - fbounds.Min.Y
	fmt.Printf("%v x %v\n", fWidth, fHeight)

	// multiply key to occupy entire new image
	// or be clever and include *factor in assignment
	// TODO: fix this mess
	for i := 0; i < fbounds.Max.X; i++ {
		for j := 0; j < fbounds.Max.Y; j++ {
			curPath := mosKeyMap[i][j]
			// open image
			curImg, err := returnImgFromPath("./input/mosaic/PCB_square_png" + "/resized" + curPath + ".png")
			if err != nil {
				fmt.Printf("Error: unable to open mosaic")
			}
			for n := 0; n < 35; n++ {
				for m := 0; m < 35; m++ {
					r, g, b, _ := curImg.At(m, n).RGBA()
					cVal := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
					finalImage.Set(i, j, cVal)
				}
			}

			// Copy current image into the new rectangle

			// copy into new img
			// save new image
		}
	}
	//fmt.Println(mosKeyMap)

	// Image Creation
	// - map mosaic images to target image pixels (distance function)
	// - create new image, write to file

	err = writeImgToFile(resizedTargetImg, "./output/resizedTarget.png")

	// TODO: Create nearest mapping function to map pixel value to nearest
	// mosaic value.

	// Profit

	fmt.Println("yipee")
}
