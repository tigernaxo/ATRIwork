package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/tigernaxo/ATRIwork/fileformat"
)

func main() {
	// 這裡是可以設定的參數
	config := newConfig(os.Args[1])
	gapRatio := config.gapRatio
	seqs := config.fileList
	refFile := config.refSeq
	outDir := config.outDir + "/"
	// intensity := config.intensity

	// 暫時先不給選
	outPNG := outDir + "histgram.png"
	histHight := 20
	convasWidth := 1920

	// 剪建立convas
	gapHeight := int(float64(histHight) * gapRatio)
	convasHeight := histHight*len(seqs) + gapHeight*(len(seqs)-1)
	convas := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{convasWidth, convasHeight},
	})

	// 讀取reference sequence
	_, refSeq := fileformat.ReadSingleFasta(refFile)
	// 預先建立給每個seq用的空siteMap、histogram array
	siteMap := make([]bool, len(refSeq))
	currentHist := make([]uint8, convasWidth)

	var startP image.Point
	// 對每個seq進行：
	for i := range seqs {
		// reset siteMap
		setAllFalse(siteMap)
		// 根據siteMap 產生currentHist
		getPointStrength(currentHist, siteMap)
		// 跳gap、繪製該histogram
		startP = image.Point{0, i * (histHight + gapHeight)}
		drawHistToConvas(currentHist, convas, color.RGBA{255, 0, 0, 255}, startP, histHight)
	}
	savePng(convas, outPNG)
}

func setAllFalse(ba []bool) {
	for i := range ba {
		ba[i] = false
	}
}

func getPointStrength(currentHist []uint8, sitMap []bool) {
	convasWidth := len(currentHist)
	unitSpan := float64(len(sitMap)) / float64(convasWidth)

	var sum float64
	var startIdx, endIdx int
	for i := 0; i < convasWidth; i++ {
		sum = 0
		switch i {
		case 0:
			startIdx = 0
			endIdx = int(unitSpan)
		case len(currentHist):
			startIdx = int(float64(convasWidth) - unitSpan)
			endIdx = convasWidth - 1
		default:
			startIdx = int(unitSpan * float64(i) / float64(convasWidth))
			endIdx = startIdx + int(unitSpan)
		}
		for j := startIdx; j < endIdx; j++ {
			if sitMap[j] {
				sum++
			}
		}
		currentHist[i] = uint8((float64(sum) / float64(endIdx-startIdx)) * float64(255))
	}
	fmt.Println(startIdx, endIdx, sum)
}

func drawHistToConvas(hist []uint8, img *image.RGBA, snvColor color.RGBA, leftTop image.Point, height int) {
	if len(hist) > img.Bounds().Max.Y {
		log.Panicf("The Histogram is out of convas bounds .\n")
	}
	xBound := int(math.Min(float64(len(hist)), float64(img.Bounds().Max.X)))
	drawArea := image.Rectangle{
		leftTop,
		image.Point{xBound, leftTop.Y + height},
	}
	for x := drawArea.Min.X; x < drawArea.Max.X; x++ {
		snvColor.R = (hist[x] / 255) * snvColor.R
		snvColor.G = (hist[x] / 255) * snvColor.G
		snvColor.B = (hist[x] / 255) * snvColor.B
		for y := drawArea.Min.Y; y < drawArea.Max.Y; y++ {
			img.SetRGBA(x, y, snvColor)
		}
	}
}

func savePng(img image.Image, outPath string) {
	outImage, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outImage.Close()
	png.Encode(outImage, img)
}
