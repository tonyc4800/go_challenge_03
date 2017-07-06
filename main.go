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

// returnImgFromPath accepts a file path to a jpeg image and returns the image.
// NOTE: `*.jpg`` is not producing expected results.
func returnImgFromPath(imgPath string) (image.Image, error) {
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open img: %v", err)
	}
	defer f.Close()

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
	rgbS := [3]uint32{0, 0, 0}
	var totalPix uint32

	// Loop image from bottom left to upper right.  Values are divided by 2^8
	// since RGBA returns values on [0, 65535](16-bit) and [0, 255](8-bit) is,
	// subjectively, easier to interpret.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgbS[0] = rgbS[0] + (r / 256)
			rgbS[1] = rgbS[1] + (g / 256)
			rgbS[2] = rgbS[2] + (b / 256)
			totalPix++
		}
	}

	// Calculate average for each channel.
	rgbS[0] = rgbS[0] / totalPix
	rgbS[1] = rgbS[1] / totalPix
	rgbS[2] = rgbS[2] / totalPix

	return rgbS
}

// resizeImage accepts and image and target width and height sizes, then resizes
// and returns the image.
func resizeImage(oImg image.Image, tWidth int, tHeight int) image.Image {

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

			// (i, j) will be the coordinates for the pix value in the new image
			// and (xStart, yStart, xCoord, yCoord) will describe the sub image.
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

	rsImgF, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to creating img file: %v", err)
	}
	defer rsImgF.Close()

	err = png.Encode(rsImgF, img)
	if err != nil {
		return fmt.Errorf("unable to write image to file: %v", err)
	}

	return nil
}

func createMosaicMapping(mosDir string, resizeMosW int, resizeMosH int) map[string][3]uint8 {

	// Create directory to hold smaller images (if not exist) 777.
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

			rsImg := resizeImage(img, resizeMosW, resizeMosH)
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

	const resizeWidth int = 120
	const resizeHeight int = 120
	const resizeMosW int = 35
	const resizeMosH int = 35

	mosMap := createMosaicMapping(mosDir, resizeMosH, resizeMosW)

	// LOOK INTO: can a map be written to a file?

	fileName := "day_man.png"
	//fileName := "boris_squat.png"

	tarImgP := "./input/target/"
	tarImgP = tarImgP + fileName

	img, err := returnImgFromPath(tarImgP)
	if err != nil {
		fmt.Printf("Error Obtaining Img: %v\n", err)
	}

	resizedTargetImg := resizeImage(img, resizeWidth, resizeHeight)

	bounds := resizedTargetImg.Bounds()
	rsWidth := bounds.Max.X - bounds.Min.X
	rsHeight := bounds.Max.Y - bounds.Min.Y

	// Loop resized image and map a mosaic value to the pixel value.
	mosKeyMap := [resizeWidth][resizeHeight]string{}
	track := 0
	for j := 0; j < rsHeight; j++ {
		for i := 0; i < rsWidth; i++ {
			r, g, b, _ := resizedTargetImg.At(i, j).RGBA()
			var mosaicN string
			closest := math.MaxFloat64
			for k, v := range mosMap {
				R := v[0]
				G := v[1]
				B := v[2]

				// Calculate nearest mosaic The squareroot is removed for
				// optimization since we don't care what the value of d is.
				rd := math.Pow((float64(R) - float64(uint8(r))), 2)
				gd := math.Pow((float64(G) - float64(uint8(g))), 2)
				bd := math.Pow((float64(B) - float64(uint8(b))), 2)

				d := rd + gd + bd
				if d < closest {
					closest = d
					mosaicN = k
				}

			}
			mosKeyMap[i][j] = mosaicN
			track++
		}
	}

	finalImage := image.NewRGBA(image.Rect(0, 0, resizeWidth*resizeMosW, resizeHeight*resizeMosH))
	//fbounds := finalImage.Bounds()
	//fWidth := fbounds.Max.X - fbounds.Min.X
	//fHeight := fbounds.Max.Y - fbounds.Min.Y

	// Loop the new mosaic image from lower left to upper right. (i, j) will be
	// used to access the resized target image. (s, t) will be used to access
	// the final image.
	s := 0
	t := 0
	for j := 0; j < resizeHeight; j++ {

		for i := 0; i < resizeWidth; i++ {
			t = resizeMosH * j

			curPath := mosKeyMap[i][j]
			curImg, err := returnImgFromPath("./input/mosaic/PCB_square_png" + "/resized/" + curPath + ".png")
			if err != nil {
				fmt.Printf("Error: unable to open mosaic: %v at [%v, %v]", curImg, i, j)
			}

			// Fill the current location with cooresponding pixel information
			// from the mosaic tile information (m,n) will be used to loop the
			// current mosaic photo.
			for n := 0; n < resizeMosH; n++ {
				s = resizeMosW * i
				for m := 0; m < resizeMosW; m++ {
					r, g, b, _ := curImg.At(m, n).RGBA()
					cVal := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
					finalImage.Set(s, t, cVal)
					s++
				}
				t++
			}

		}

		fmt.Printf("col complete: %v of %v\n", j, resizeHeight)
	}

	err = writeImgToFile(resizedTargetImg, "./output/resizedTarget.png")
	if err != nil {
		fmt.Println("Error writing the resized target image to file")
	}

	outPath := "./output/" + fileName
	err = writeImgToFile(finalImage, outPath)
	if err != nil {
		fmt.Println("Error writing the final mosaic image to file")
	}

	fmt.Println("yipee")
}
