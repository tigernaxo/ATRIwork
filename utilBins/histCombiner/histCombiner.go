package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/nfnt/resize"
)

// Arg[1]:gap/histogram ratio。
// Arg[2]:folder，讀取folder/list裏面的list file進行combine。
func main() {
	// Setting arguments to variables
	ratio, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		panic(err)
	}
	folder := os.Args[2] + "/"
	metaName := "list"
	outName := "combinedImage.png"

	// Opening forder list
	fileList := []string{}
	f, err := os.Open(folder + metaName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fileList = append(fileList, folder+scanner.Text())
	}
	log.Println(fileList)

	// Calculate convase Bounds, then create
	histHeight := 20
	gapHeight := int(float64(histHeight) * ratio)
	convasHeight := histHeight*len(fileList) + len(fileList) - 1*gapHeight

	imgData, _ := decodeImg(fileList[0])
	convasWidth := imgData.Bounds().Max.X

	convas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{convasWidth, convasHeight},
	})

	// Draw images and gap
	finishedHeight := 0
	for _, file := range fileList {
		fmt.Println(file)
		img, _ := decodeImg(file)
		if !checkWidth(img, convasWidth) {
			panic("the image files do not have equal width")
		}
		startY := finishedHeight
		endY := finishedHeight + histHeight - 1
		for x := 0; x < convasWidth; x++ {
			// Get image color, and set convase color
			_, g, b, a := img.At(x, 0).RGBA()
			pointColor := color.RGBA{
				uint8(255),
				uint8(g),
				uint8(b),
				uint8(a),
			}
			// Drawing histogram on convas
			for y := startY; y <= endY; y++ {
				convas.SetRGBA(x, y, pointColor)
			}
		}
		// Skip gap height
		fmt.Println("convase X, Y:", convasWidth, convasHeight)
		fmt.Println("StartY, EndY: ", startY, endY)
		finishedHeight = finishedHeight + histHeight + gapHeight
	}
	combinedImage := resize.Resize(
		uint(1920),
		uint(1080),
		convas,
		resize.NearestNeighbor,
	)

	// Create png file to output
	outImage, err := os.Create(folder + outName)
	if err != nil {
		log.Fatal(err)
	}
	defer outImage.Close()

	png.Encode(outImage, combinedImage)
}

func checkWidth(img image.Image, width int) bool {
	switch img.Bounds().Max.X {
	case width:
		return true
	default:
		return false
	}
}

func decodeImg(fileName string) (image.Image, string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	imgData, imgType, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return imgData, imgType
}
