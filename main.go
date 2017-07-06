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

// calcAvgRGB accepts and image and returns the average pixel values for each
// channel as an three value array (RGB).
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
func resizeImage(oImg image.Image, tW int, tH int) image.Image {

	// Create new, resized, image rectangle and obtain the ratio need to
	// scale the original image down into the target image.  The ratio is used
	// to define subimage pixel bounds on the original image.
	rsImg := image.NewRGBA(image.Rect(0, 0, tW, tH))
	bounds := oImg.Bounds()
	oW := bounds.Max.X - bounds.Min.X
	oH := bounds.Max.Y - bounds.Min.Y
	wRatio := float64(oW) / float64(tW)
	hRatio := float64(oH) / float64(tH)

	// Create a grid of upper right/bounding coordinates for subimages from the
	// original image that can be mapped into the new image. Coordinate values
	// will be cropped to an int value, not rounded.
	var yBoundsU []int
	var xBoundsU []int

	// Do not loop until the target height and width since we are multiplying
	// the upper bound by the ratio.  The original max value will be added to
	// appended the the slice after looping.
	for y := 0; y < tH; y++ {
		i := int(float64(y) * hRatio)
		yBoundsU = append(yBoundsU, i)
	}
	yBoundsU = append(yBoundsU, bounds.Max.Y)

	for x := 0; x < tW; x++ {
		i := int(float64(x) * wRatio)
		xBoundsU = append(xBoundsU, i)
	}
	xBoundsU = append(xBoundsU, bounds.Max.X)

	// Remove first value from slice since it is not an upper bound.
	xBoundsU = append(xBoundsU[:0], xBoundsU[0+1:]...)
	yBoundsU = append(yBoundsU[:0], yBoundsU[0+1:]...)

	// Create sub images from the original image.  The subimage bounds are
	// contained by a rectangle defined as (xBoundL, yBoundL, xBoundU, yBoundU).

	yBoundL := 0
	for j, yBoundU := range yBoundsU {
		xBoundL := 0
		for i, xBoundU := range xBoundsU {

			// (i, j) will be the coordinates for the pix value in the new image
			// and (xBoundL, yBoundL, xBoundU, yBoundU) will describe the sub image.
			subImg := image.NewRGBA(image.Rect(0, 0, xBoundU-xBoundL, yBoundU-yBoundL))

			// Fill subimage pixel values.
			n := 0
			for yy := yBoundL; yy <= yBoundU; yy++ {
				m := 0
				for xx := xBoundL; xx <= xBoundU; xx++ {
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
			rsImg.Set(i, j, nVal)
			xBoundL = xBoundU
		}
		yBoundL = yBoundU
	}

	return rsImg
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

func createMosaicMapping(mosDir string, rsMosW int, rsMosH int) map[string][3]uint8 {

	// Create directory to hold smaller images (if not exist) os.ModePerm is
	// equivalent to unix permissions `777`.
	rsDir := mosDir + "/resized"
	if _, err := os.Stat(rsDir); os.IsNotExist(err) {
		os.Mkdir(rsDir, os.ModePerm)
	}

	mosMap := make(map[string][3]uint8)
	mosFiles, _ := ioutil.ReadDir(mosDir)

	for _, f := range mosFiles {
		fPath := f.Name()
		ext := filepath.Ext(fPath)
		key := fPath[0 : len(fPath)-len(ext)]

		if ext == ".png" {
			img, err := returnImgFromPath(mosDir + "/" + fPath)
			if err != nil {
				fmt.Println(fPath)
				fmt.Printf("Error Obtaining Img: %v\n", err)
			}

			rsImg := resizeImage(img, rsMosW, rsMosH)
			rsPath := rsDir + "/" + fPath
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

	const rsTarW int = 300
	const rsTarH int = 300
	const rsMosW int = 35
	const rsMosH int = 35

	mosMap := createMosaicMapping(mosDir, rsMosH, rsMosW)

	tarImgDir := "./input/target/"
	//tarName := "day_man.png"
	//tarName := "boris_squat.png"
	tarName := "abby_jack.png"
	tarImgP := tarImgDir + tarName

	img, err := returnImgFromPath(tarImgP)
	if err != nil {
		fmt.Printf("Error Obtaining Img: %v\n", err)
	}

	rsTarImg := resizeImage(img, rsTarW, rsTarH)

	// Loop resized image and map a mosaic value to the pixel value.
	mosKeyMap := [rsTarW][rsTarH]string{}
	bounds := rsTarImg.Bounds()

	// Loop resized image columns (height).
	for j := 0; j < (bounds.Max.Y - bounds.Min.Y); j++ {

		// Loop resized image columns (width).
		for i := 0; i < (bounds.Max.X - bounds.Min.X); i++ {
			r, g, b, _ := rsTarImg.At(i, j).RGBA()
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
		}
	}

	finalImg := image.NewRGBA(image.Rect(0, 0, rsTarW*rsMosW, rsTarH*rsMosH))

	// Loop the new mosaic image from lower left to upper right. (i, j) will be
	// used to access the resized target image. (s, t) will be used to access
	// the final image.
	s := 0
	t := 0
	for j := 0; j < rsTarH; j++ {
		for i := 0; i < rsTarW; i++ {
			t = rsMosH * j
			curPath := mosKeyMap[i][j]
			curImg, err := returnImgFromPath("./input/mosaic/PCB_square_png" + "/resized/" + curPath + ".png")
			if err != nil {
				fmt.Printf("Error: unable to open mosaic: %v at [%v, %v]", curImg, i, j)
			}

			// Fill the current location with cooresponding pixel information
			// from the mosaic tile information (m,n) will be used to loop the
			// current mosaic photo.
			for n := 0; n < rsMosH; n++ {
				s = rsMosW * i
				for m := 0; m < rsMosW; m++ {
					r, g, b, _ := curImg.At(m, n).RGBA()
					cVal := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
					finalImg.Set(s, t, cVal)
					s++
				}
				t++
			}
		}
		fmt.Printf("col complete: %v of %v\n", j+1, rsTarH)
	}

	err = writeImgToFile(rsTarImg, "./output/resizedTarget.png")
	if err != nil {
		fmt.Println("Error writing the resized target image to file")
	}

	outPath := "./output/" + tarName
	err = writeImgToFile(finalImg, outPath)
	if err != nil {
		fmt.Println("Error writing the final mosaic image to file")
	}

	fmt.Println("yipee")
}
